# Stage 1: Build the Go binary
FROM golang:1.23 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go Modules files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod tidy

# Copy the source code into the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Stage 2: Create the final minimal image
FROM alpine:latest

# Install ca-certificates to allow connecting to secure sites
RUN apk --no-cache add ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /app

USER 1000

# Copy the Go binary from the builder stage
COPY --from=builder /app/main .

# Run the Go binary
CMD ["./rss-collector"]