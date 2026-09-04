#!/usr/bin/env bash
#
# Wraps `influxdb3 serve` so an unclean shutdown (host crash, power loss, `docker kill`)
# that leaves a single WAL segment or dbs/ file half-written doesn't turn into a permanent
# crash-loop. On each crash it greps the run's own output for the file the error names,
# and if that file is safely identifiable (under wal/ or dbs/, inside the data dir), deletes
# just that file and retries — mirrors the manual fix in
# docs/eh-services/install-influxdb.md#corrupted-file-causing-crash-loop--repeated-errors-delete-just-that-file
#
# Anything else (catalog/ corruption, no recognizable path, repeated failures past the
# attempt cap) is left to crash-loop as before so it stays visible — this never touches
# the admin token/catalog and never wipes data on its own.

set -uo pipefail

MAX_HEAL_ATTEMPTS="${INFLUXDB_HEAL_MAX_ATTEMPTS:-3}"

data_dir=""
node_id=""
args=("$@")
for i in "${!args[@]}"; do
  case "${args[$i]}" in
    --data-dir=*) data_dir="${args[$i]#--data-dir=}" ;;
    --data-dir) data_dir="${args[$((i + 1))]:-}" ;;
    --node-id=*) node_id="${args[$i]#--node-id=}" ;;
    --node-id) node_id="${args[$((i + 1))]:-}" ;;
  esac
done
data_dir="${data_dir:-${INFLUXDB3_DATA_DIR:-/var/lib/influxdb3/data}}"
data_dir_resolved="$(readlink -f -- "$data_dir")"

# influxdb3 stores everything under data-dir/<node-id>/... (its "prefix") and names files in
# error output the same way, e.g. `path=ehm-influxdb/wal/00000024877.wal` — so a bare
# `data-dir/wal/...` guess misses. Try the node-id-prefixed location first, then the bare one.
candidate_dirs=("$data_dir")
[[ -n "$node_id" ]] && candidate_dirs=("${data_dir}/${node_id}" "$data_dir")

log() { echo "[influxdb-heal] $*"; }

child_pid=""
forward_signal() {
  [[ -n "$child_pid" ]] && kill -TERM "$child_pid" 2>/dev/null
}
trap forward_signal TERM INT

attempt=0
while :; do
  attempt=$((attempt + 1))
  out="$(mktemp)"
  log "starting influxdb3 (attempt ${attempt}/$((MAX_HEAL_ATTEMPTS + 1)))"

  (
    exec > >(tee -a "$out") 2>&1
    exec /usr/bin/entrypoint.sh "$@"
  ) &
  child_pid=$!
  wait "$child_pid"
  status=$?
  child_pid=""
  sleep 0.3 # let the tee'd copy in $out catch up with the last lines before we grep it

  if [[ "$status" -eq 0 ]]; then
    rm -f "$out"
    exit 0
  fi

  if [[ "$attempt" -gt "$MAX_HEAL_ATTEMPTS" ]]; then
    log "gave up after ${MAX_HEAL_ATTEMPTS} auto-heal attempts (exit ${status}) — leaving it crashed for manual recovery, see docs/eh-services/install-influxdb.md"
    rm -f "$out"
    exit "$status"
  fi

  file_rel="$(grep -iE 'failed to (read|parse|replay)' "$out" | grep -oE '(wal|dbs)/[A-Za-z0-9_./-]+' | tail -n1)"
  if [[ -z "$file_rel" ]]; then
    file_rel="$(grep -iE 'error|corrupt|invalid' "$out" | grep -oE '(wal|dbs)/[A-Za-z0-9_./-]+' | tail -n1)"
  fi
  rm -f "$out"

  if [[ -z "$file_rel" ]]; then
    log "crashed (exit ${status}) but no wal/ or dbs/ file was named in the output — not auto-healing, exiting so the crash-loop is visible"
    exit "$status"
  fi

  case "$file_rel" in
    wal/* | dbs/*) ;;
    *)
      log "offending path '${file_rel}' is outside wal/ and dbs/ (likely catalog/) — that isn't safely fixable file-by-file, exiting for manual recovery"
      exit "$status"
      ;;
  esac

  resolved=""
  for dir in "${candidate_dirs[@]}"; do
    target="${dir}/${file_rel}"
    candidate_resolved="$(readlink -f -- "$target" 2>/dev/null || true)"

    if [[ -z "$candidate_resolved" || "$candidate_resolved" != "$data_dir_resolved"/* ]]; then
      log "refusing to consider '${target}' — resolves outside the data directory"
      continue
    fi

    if [[ -e "$candidate_resolved" ]]; then
      resolved="$candidate_resolved"
      break
    fi
  done

  if [[ -z "$resolved" ]]; then
    tried=""
    for dir in "${candidate_dirs[@]}"; do tried="${tried}${dir}/${file_rel}, "; done
    log "output named '${file_rel}' but no matching file was found on disk (tried: ${tried%, }) — nothing to delete"
    exit "$status"
  fi

  log "deleting corrupted file named in the crash output: ${resolved}"
  rm -f -- "$resolved"
done
