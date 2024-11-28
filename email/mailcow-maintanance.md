# Backup
Ref (General Backup): https://docs.mailcow.email/backup_restore/b_n_r-backup/
Ref (MySqlDump): https://docs.mailcow.email/backup_restore/b_n_r-backup_restore-mysql/
- Create backup directoris: 
```bash
mkdir -p /opt/backups/mailcow/all
mkdir -p /opt/backups/mailcow/mysqldump
```
- Take backup: 
```bash
# Manual backup (All Services):
MAILCOW_BACKUP_LOCATION=/opt/backups/mailcow/all THREADS=14 /root/mailcow-dockerized/helper-scripts/backup_and_restore.sh backup all

# manual mysqldump :
```bash
cd /root/mailcow-dockerized
source mailcow.conf
DATE=$(date +"%Y%m%d_%H%M%S")
docker compose exec -T mysql-mailcow mysqldump --default-character-set=utf8mb4 -u${DBUSER} -p${DBPASS} ${DBNAME} > /opt/backups/mailcow/mysqldump/backup_${DBNAME}_${DATE}.sql
```

# Auto backup (All Services):
https://docs.mailcow.email/backup_restore/b_n_r-backup/#cronjob

# Auto mysqldump:
Use backup script for mysqldump in cronejob