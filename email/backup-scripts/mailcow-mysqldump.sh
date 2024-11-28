#!/bin/bash

# Define absolute paths
MAILCOW_DIR="/root/mailcow-dockerized"
CONFIG_FILE="${MAILCOW_DIR}/mailcow.conf"
BACKUP_DIR="/opt/backups/mailcow/mysqldump"
DATE=$(date +"%Y%m%d_%H%M%S")

# Ensure backup directory exists
mkdir -p "$BACKUP_DIR"

# Load Mailcow configuration
if [ -f "$CONFIG_FILE" ]; then
    source "$CONFIG_FILE"
else
    echo "Configuration file $CONFIG_FILE not found!"
    exit 1
fi

# Perform MySQL dump
if [ -z "$DBUSER" ] || [ -z "$DBPASS" ] || [ -z "$DBNAME" ]; then
    echo "Database configuration variables are not set!"
    exit 1
fi

docker compose -f "${MAILCOW_DIR}/docker-compose.yml" exec -T mysql-mailcow mysqldump \
    --default-character-set=utf8mb4 \
    -u"$DBUSER" \
    -p"$DBPASS" \
    "$DBNAME" > "${BACKUP_DIR}/backup_${DBNAME}_${DATE}.sql"

# Ensure success
if [ $? -eq 0 ]; then
    echo "Backup completed successfully: ${BACKUP_DIR}/backup_${DBNAME}_${DATE}.sql"
else
    echo "Backup failed!"
    exit 1
fi
