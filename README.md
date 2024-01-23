#  Golang: How to deploy WbAPI CRUD for MongoDb to Docker Desktop

We are going to **containerize** our **Go application** and **MongoDB** in two Docker containers running on our local laptop, we'll need to create a **Dockerfile** for our Go application

You'll also need to ensure that **both containers can communicate with each other**, usually done through **Docker networking** features like user-defined networks or linking

Below is an example Dockerfile for your Go application and additional steps to set up the environment properly

## 1. Dockerfile for Go Application

We create the Dockerfile and we save it in the same directory as our Go application source code

**Dockerfile**

```
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
```

## 2. Build Go Application Docker Image

```
docker build -t my-go-app .
```

## 3. Network Configuration

Since your MongoDB container is already running, you'll need to ensure your Go application can connect to it. 

If your MongoDB container is running with the default settings, it should be accessible via localhost on your host machine, but from another container, you need to use Docker networking.

## 4. Create a Docker Network (if you haven't already)

```
docker network create my-network
```

## 5. Connect MongoDB Container to our Network

Assuming your MongoDB container is named my-mongo, run:

```
docker network connect my-network mongodb
```

## 6. Run Your Go Application Container:

When running your Go application container, you should also attach it to the same network

You'll need to adjust the MongoDB URI in your Go application to use the name of the MongoDB container as the hostname, e.g., mongodb://mongodb:2701

```
docker run -d --name my-go-app-instance --network my-network -p 8080:8080 my-go-app
```


This setup uses Docker's default bridge network to facilitate communication between the two containers using the container names as hostnames for inter-container communication. Ensure your application's MongoDB URI reflects this network configuration.
