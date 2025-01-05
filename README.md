# RSS Collector

This is a project that I made for myself. This project contains three micro services that consist one service.
Why over complicated? Because it is a hobby/learning project. Feel free to use it, of course.

Simplest way to find them on Dockerhub and make your own docker-compose based on [this](https://github.com/onlyati/rss-collector/blob/main/docker-compose.yaml) file.
- https://hub.docker.com/r/onlyati/rss-collector
- https://hub.docker.com/r/onlyati/rss-api
- https://hub.docker.com/r/onlyati/rss-processor

## What is this?

I have made this project to collect some RSS feed for myself, but with a twist:
```
  .-----------.      .-------.      .-----------.      .------------.      .-----------.
  | Collector | ---> | Kafka | ---> | Processor | ---> | PostgreSQL | ---> | REST  API |
  |  service  |      '-------'      |  service  |      '------------'      |  service  |
  '-----------'                     '-----------'                          '-----------'
```

At the end this service run in my kubernetes cluster, the processor and REST service can be scaled up.

## Environment variables for container

Following environment variables can be used in all container:
- `RSS_CONFIG_PATH`: If you specify a path and mount a config map there, it would be used.

REST API service variables:
- `RSS_DB_HOSTNAME`: PostgreSQL database address
- `RSS_DB_PORT`: PostgreSQL port number
- `RSS_DB_NAME`: Database name
- `RSS_DB_USER`: User for the database
- `RSS_DB_PW_PATH`: Path for a file which store the password (file can be mounted via secret)
- `RSS_API_HOSTNAME`: Hostname where REST API listen
- `RSS_API_PORT`: Port number where REST API listen
- `RSS_API_AUTH_CONFIG`: Address for keycloak endpoint list, example https://{{Hostname}}/realms/{{Realm}}/.well-known/openid-configuration
- `RSS_CORS_ORIGINS`: Origins list, separated by ',' for allowed origins
- `RSS_CORS_METHODS`: Method list, separated by ',' for allowed methods

Collector service variables:
- `RSS_YOUTUBE_*`: List about Youtube channels that must be collected
- `RSS_REDDIT_*`: List about Reddit threads that must be collected
- `RSS_STANDARD_*`: Any RSS feed
- `RSS_CRUNCHYROLL`: If it has value then Crunchyroll new anime feed is collected
- `RSS_KAFKA_SERVER`: Kafka address
- `RSS_KAFKA_TOPIC`: Kafka topic
- `RSS_WAIT_SECONDS`: Wait time between two collections

Processor service variables:
- `RSS_DB_HOSTNAME`: PostgreSQL database address
- `RSS_DB_PORT`: PostgreSQL port number
- `RSS_DB_NAME`: Database name
- `RSS_DB_USER`: User for the database
- `RSS_DB_PW_PATH`: Path for a file which store the password (file can be mounted via secret)
- `RSS_KAFKA_SERVER`: Kafka address
- `RSS_KAFKA_TOPIC`: Kafka topic
- `RSS_KAFKA_GROUP_ID`: Kafka consumer group is
