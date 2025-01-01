#!/bin/bash

# Backup script for both monolithic and microservices mode

# Load environment variables
set -a
source .env
set +a

# Create backup directory
timestamp=$(date +%Y%m%d_%H%M%S)
backup_dir="./backups/$timestamp"
mkdir -p $backup_dir

# Function to backup a database
backup_database() {
    local host=$1
    local port=$2
    local user=$3
    local password=$4
    local db=$5
    local output=$6
    
    PGPASSWORD=$password pg_dump -h $host -p $port -U $user -d $db -F c > "$output"
}

# Function to backup Redis
backup_redis() {
    local host=$1
    local port=$2
    local output=$3
    
    redis-cli -h $host -p $port --rdb "$output"
}

if [ "$API_MODE" = "monolithic" ]; then
    echo "Backing up monolithic mode..."
    
    # Backup PostgreSQL databases
    backup_database "localhost" "5432" "$DB_USER" "$DB_PASSWORD" "users" "$backup_dir/users.dump"
    backup_database "localhost" "5432" "$DB_USER" "$DB_PASSWORD" "notifications" "$backup_dir/notifications.dump"
    
    # Backup Redis
    backup_redis "localhost" "6379" "$backup_dir/redis.rdb"
    
    # Backup Elasticsearch
    curl -X GET "http://localhost:9200/_snapshot/my_backup" -H 'Content-Type: application/json' -d'
    {
        "type": "fs",
        "settings": {
            "location": "'$backup_dir'/elasticsearch"
        }
    }'
    
else
    echo "Backing up microservices mode..."
    
    # Backup User Service database
    backup_database "$USER_DB_HOST" "$USER_DB_PORT" "$USER_DB_USER" "$USER_DB_PASSWORD" "$USER_DB_NAME" "$backup_dir/users.dump"
    backup_redis "$USER_REDIS_HOST" "$USER_REDIS_PORT" "$backup_dir/user_redis.rdb"
    
    # Backup Notification Service database
    backup_database "$NOTIFICATION_DB_HOST" "$NOTIFICATION_DB_PORT" "$NOTIFICATION_DB_USER" "$NOTIFICATION_DB_PASSWORD" "$NOTIFICATION_DB_NAME" "$backup_dir/notifications.dump"
    backup_redis "$NOTIFICATION_REDIS_HOST" "$NOTIFICATION_REDIS_PORT" "$backup_dir/notification_redis.rdb"
    
    # Backup Consul data
    consul snapshot save "$backup_dir/consul.snap"
fi

# Compress backup
tar -czf "$backup_dir.tar.gz" -C "./backups" "$timestamp"
rm -rf $backup_dir

echo "Backup completed successfully!"
echo "Backup file: $backup_dir.tar.gz" 