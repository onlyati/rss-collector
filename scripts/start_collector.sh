#!/bin/sh

# 
# If config path is specified, then use it
#
if [ "$RSS_CONFIG_PATH" != "" ]; then
    echo "start application with ${RSS_CONFIG_PATH} configuration"
    exec /app/rss-collector collect --config "$RSS_CONFIG_PATH"
fi

#
# Else collect configuration from environment variables
# RSS_YOUTUBE_*   => Youtube channels
# RSS_REDDIT_*    => Reddit threads
# RSS_STANDARD_*  => General RSS feeds
# RSS_CRUNCHYROLL => If value specified then anime news feed is read
#
echo "collect config data from environment variables"
echo "collect RSS_YOUTUBE_ variables"
youtube=""
for line in $(env | grep -e "^RSS_YOUTUBE_"); do
    channel=$(echo "$line" | cut -d'=' -f2)
    youtube="$youtube $channel"
done

echo "collect RSS_REDDIT_ variables"
reddit=""
for line in $(env | grep -e "^RSS_REDDIT_"); do
    channel=$(echo "$line" | cut -d'=' -f2)
    reddit="$reddit $channel"
done

echo "collect RSS_STANDARD_ variables"
standard=""
for line in $(env | grep -e "^RSS_STANDARD_"); do
    channel=$(echo "$line" | cut -d'=' -f2)
    standard="$standard $channel"
done

echo "checking value of RSS_CRUNCHYROLL"
crunchyroll="$RSS_CRUNCHYROLL"

#
# Generate config file
#
echo "generating configuration"
file=/app/config.yaml

echo "youtube:" > "$file"
for elem in $youtube; do
    echo "- $elem" >> "$file"
done

echo "reddit:" >> "$file"
for elem in $reddit; do
    echo "- $elem" >> "$file"
done

echo "standard:" >> "$file"
for elem in $standard; do
    echo "- $elem" >> "$file"
done

if [ "$crunchyroll" != "" ]; then
    echo "crunchyroll: true" >> "$file"
fi

#
# Add RSS_KAFKA_ settings
#
echo "kafka:" >> "$file"
echo "  server: ${RSS_KAFKA_SERVER}" >> "$file"
echo "  topic: ${RSS_KAFKA_TOPIC}" >> "$file"
echo "wait_seconds: ${RSS_WAIT_SECONDS}" >> "$file"

echo "------ GENERATED CONFIG ----------"
cat "$file"
echo "------ GENERATED CONFIG ----------"

exec /app/rss-collector collect --config "$file"
