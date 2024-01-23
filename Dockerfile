# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
FROM golang:1.21.6 as builder

# Copy local code to the container image.
WORKDIR /app
COPY . .

# Download Go modules and build your application.
# Consider using go mod tidy before building to ensure all dependencies are correct.
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -v -o server

# Use a Docker multi-stage build to create a lean production image.
FROM alpine:latest  
RUN apk --no-cache add ca-certificates

WORKDIR /root/
# Copy the binary from the builder stage to the production image.
COPY --from=builder /app/server .

# Expose the port your app runs on.
EXPOSE 8080

# Run the binary.
CMD ["./server"]
