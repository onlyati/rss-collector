#!/bin/sh

# 
# If config path is specified, then use it
#
if [ "$RSS_CONFIG_PATH" != "" ]; then
    echo "start application with ${RSS_CONFIG_PATH} configuration"
    exec /app/rss-api listen --config "$RSS_CONFIG_PATH"
fi

#
# Else read config from environment variables
#
file=/app/config.yaml
echo "db:" > $file
echo "  hostname: $RSS_DB_HOSTNAME" >> $file
echo "  port: $RSS_DB_PORT" >> $file
echo "  user: $RSS_DB_USER" >> $file
echo "  password_path: $RSS_DB_PW_PATH" >> $file
echo "  db_name: $RSS_DB_NAME" >> $file
echo "api:" >> $file
echo "  hostname: $RSS_API_HOSTNAME" >> $file
echo "  port: $RSS_API_PORT" >> $file

echo "------ GENERATED CONFIG ----------"
cat "$file"
echo "------ GENERATED CONFIG ----------"

# Start api
exec /app/rss-api listen --config "$file"
