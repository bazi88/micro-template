#!/bin/bash

# Restore script for both monolithic and microservices mode

# Load environment variables
set -a
source .env
set +a

# Check if backup file is provided
if [ -z "$1" ]; then
    echo "Please provide backup file path"
    echo "Usage: $0 <backup_file>"
    exit 1
fi

backup_file=$1
temp_dir="/tmp/micro_restore_$(date +%s)"

# Extract backup
mkdir -p $temp_dir
tar -xzf $backup_file -C $temp_dir

# Function to restore a database
restore_database() {
    local host=$1
    local port=$2
    local user=$3
    local password=$4
    local db=$5
    local input=$6
    
    PGPASSWORD=$password pg_restore -h $host -p $port -U $user -d $db -c $input
}

# Function to restore Redis
restore_redis() {
    local host=$1
    local port=$2
    local input=$3
    
    redis-cli -h $host -p $port FLUSHALL
    redis-cli -h $host -p $port --rdb $input
}

if [ "$API_MODE" = "monolithic" ]; then
    echo "Restoring monolithic mode..."
    
    # Stop all services
    docker-compose -f docker-compose.monolithic.yml down
    
    # Restore PostgreSQL databases
    restore_database "localhost" "5432" "$DB_USER" "$DB_PASSWORD" "users" "$temp_dir/users.dump"
    restore_database "localhost" "5432" "$DB_USER" "$DB_PASSWORD" "notifications" "$temp_dir/notifications.dump"
    
    # Restore Redis
    restore_redis "localhost" "6379" "$temp_dir/redis.rdb"
    
    # Restore Elasticsearch
    curl -X POST "http://localhost:9200/_snapshot/my_backup/_restore" -H 'Content-Type: application/json' -d'
    {
        "indices": "*",
        "include_global_state": true
    }'
    
    # Start services
    docker-compose -f docker-compose.monolithic.yml up -d
    
else
    echo "Restoring microservices mode..."
    
    # Stop all services
    docker-compose -f docker-compose.microservices.yml down
    
    # Restore User Service
    ssh $USER_SERVICE_HOST "cd ~/micro/services/user && docker-compose down"
    restore_database "$USER_DB_HOST" "$USER_DB_PORT" "$USER_DB_USER" "$USER_DB_PASSWORD" "$USER_DB_NAME" "$temp_dir/users.dump"
    restore_redis "$USER_REDIS_HOST" "$USER_REDIS_PORT" "$temp_dir/user_redis.rdb"
    ssh $USER_SERVICE_HOST "cd ~/micro/services/user && docker-compose up -d"
    
    # Restore Notification Service
    ssh $NOTIFICATION_SERVICE_HOST "cd ~/micro/services/notification && docker-compose down"
    restore_database "$NOTIFICATION_DB_HOST" "$NOTIFICATION_DB_PORT" "$NOTIFICATION_DB_USER" "$NOTIFICATION_DB_PASSWORD" "$NOTIFICATION_DB_NAME" "$temp_dir/notifications.dump"
    restore_redis "$NOTIFICATION_REDIS_HOST" "$NOTIFICATION_REDIS_PORT" "$temp_dir/notification_redis.rdb"
    ssh $NOTIFICATION_SERVICE_HOST "cd ~/micro/services/notification && docker-compose up -d"
    
    # Restore Consul data
    consul snapshot restore "$temp_dir/consul.snap"
    
    # Start services
    docker-compose -f docker-compose.microservices.yml up -d
fi

# Cleanup
rm -rf $temp_dir

echo "Restore completed successfully!" 