#!/bin/sh

# 
# If config path is specified, then use it
#
if [ "$RSS_CONFIG_PATH" != "" ]; then
    # Check if need database migration
    if [ "$RSS_DB_MIGRATE" != "" ]; then
        echo "start database migration"
        /app/rss-processor db-migration --config "$RSS_CONFIG_PATH" || exit 20
    fi
    echo "start application with ${RSS_CONFIG_PATH} configuration"
    exec /app/rss-processor process --config "$RSS_CONFIG_PATH"
fi

#
# Else read config from environment variables
#
file=/app/config.yaml
echo "generating configuration"
echo "kafka:" > $file
echo "  server: $RSS_KAFKA_SERVER" >> $file
echo "  topic: $RSS_KAFKA_TOPIC" >> $file
echo "  group_id: $RSS_KAFKA_GROUP_ID" >> $file
echo "db:" >> $file
echo "  hostname: $RSS_DB_HOSTNAME" >> $file
echo "  port: $RSS_DB_PORT" >> $file
echo "  user: $RSS_DB_USER" >> $file
echo "  password_path: $RSS_DB_PW_PATH" >> $file
echo "  db_name: $RSS_DB_NAME" >> $file

echo "------ GENERATED CONFIG ----------"
cat "$file"
echo "------ GENERATED CONFIG ----------"

# Check if need database migration
if [ "$RSS_DB_MIGRATE" != "" ]; then
    echo "start database migration"
    /app/rss-processor db-migration --config "$file" || exit 20
fi
exec /app/rss-processor process --config "$file"
