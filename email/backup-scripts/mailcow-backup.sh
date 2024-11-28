#!/bin/sh

# Backup mailcow data
# https://docs.mailcow.email/backup_restore/b_n_r-backup/

set -e

OUT="$(mktemp)"
export MAILCOW_BACKUP_LOCATION="/opt/backups/mailcow/all"
SCRIPT="/root/mailcow-dockerized/helper-scripts/backup_and_restore.sh"
PARAMETERS="backup all"
OPTIONS="--delete-days 30"

# run command
set +e
"${SCRIPT}" ${PARAMETERS} ${OPTIONS} 2>&1 > "$OUT"
RESULT=$?

if [ $RESULT -ne 0 ]
    then
            echo "${SCRIPT} ${PARAMETERS} ${OPTIONS} encounters an error:"
            echo "RESULT=$RESULT"
            echo "STDOUT / STDERR:"
            cat "$OUT"
fi