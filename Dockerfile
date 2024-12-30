# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder

RUN apk --no-cache add ca-certificates librdkafka-dev pkgconf musl-dev alpine-sdk
# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go Modules files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod tidy

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -tags musl -o rss-collector .

# Stage 2: Create the final minimal image
FROM alpine:latest

# Install ca-certificates to allow connecting to secure sites
RUN apk --no-cache add ca-certificates librdkafka-dev pkgconf

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go binary from the builder stage
COPY --from=builder /app/rss-collector .
COPY ./test_source/test_config.yaml /app/config.yaml

# Run the Go binary
ENTRYPOINT ["/app/rss-collector", "collect", "--config", "/app/config.yaml"]

