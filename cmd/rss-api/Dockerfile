# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go Modules files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod tidy

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -tags musl -o rss-api ./cmd/rss-api/main.go

# Stage 2: Create the final minimal image
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go binary from the builder stage
COPY --from=builder /app/rss-api .
COPY ./scripts/start_api.sh /app/start_api.sh
COPY ./openapi /app/openapi

# Install ca-certificates to allow connecting to secure sites
RUN addgroup -S rss \
    && adduser -S rss -G rss \
    && apk --no-cache add ca-certificates pkgconf \
    && chmod +x /app/start_api.sh \
    && chown -R rss:rss /app

USER rss
# Run the Go binary
ENTRYPOINT ["/app/start_api.sh"]

