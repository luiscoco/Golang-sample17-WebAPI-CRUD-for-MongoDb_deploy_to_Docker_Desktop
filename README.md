#  Golang: How to deploy WebAPI CRUD for MongoDb to Docker Desktop

We are going to **containerize** our **Go application** and **MongoDB** in two Docker containers running on our local laptop, we'll need to create a **Dockerfile** for our Go application

You'll also need to ensure that **both containers can communicate with each other**, usually done through **Docker networking** features like user-defined networks or linking

Below is an example **Dockerfile for your Go application** and additional steps to set up the environment properly


## 1. Pull and run MongoDB docker container

We first pull the mongodb image

```
docker pull mongo
```

We run the mongo docker container

```
docker run --name mongodb -d -p 27017:27017 --restart unless-stopped mongo
```

![image](https://github.com/luiscoco/Golang-sample16-WebAPI-CRUD-for-MongoDb/assets/32194879/60d27a6f-edbb-4116-90c3-3ac8346fd813)

- Verify the image and container in **Docker Desktop**

![image](https://github.com/luiscoco/Golang-sample16-WebAPI-CRUD-for-MongoDb/assets/32194879/5a959223-0fbe-46d8-be07-6d2136f99807)

![image](https://github.com/luiscoco/Golang-sample16-WebAPI-CRUD-for-MongoDb/assets/32194879/cda014ad-a77c-4fd1-a96b-2ab4770bbf12)

## 2. Dockerfile for Go Application

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

## 3. Build Go Application Docker Image

```
docker build -t my-go-app .
```

## 4. Network Configuration

Since your MongoDB container is already running, we'll need to ensure our Go application can connect to it

If our MongoDB container is running with the default settings, it should be accessible via localhost on your host machine, but from another container, we need to use Docker networking

## 5. Create a Docker Network

```
docker network create my-network
```

## 6. Connect MongoDB Container to our Network

Assuming our MongoDB container is named "mongodb", run:

```
docker network connect my-network mongodb
```

## 7. Run Your Go Application Container

When running our Go application container, we should also attach it to the same network

We'll need to adjust the MongoDB URI in our Go application to use the name of the MongoDB container as the hostname, e.g., mongodb://mongodb:27017

**IMPORTANT NOTE**: we have to modify the following line in the **main.go** file to confirm the mongo container hostname

clientOptions := options.Client().ApplyURI("**mongodb://mongodb:27017**")

```
func connectDB() {
    // Change the URI to "mongodb://mongodb:27017" to connect to the MongoDB container
    clientOptions := options.Client().ApplyURI("mongodb://mongodb:27017")
    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    collection = client.Database("bookstore").Collection("books")
}
```

We also run the following command to execute the Go application docker container

```
docker run -d --name my-go-app-instance --network my-network -p 8080:8080 my-go-app
```

We see the docker images and container in Docker Desktop

![image](https://github.com/luiscoco/Golang-sample17-WebAPI-CRUD-for-MongoDb_deploy_to_Docker_Desktop/assets/32194879/1628d0c3-df01-48e3-b5f5-5271c07dca7e)

![image](https://github.com/luiscoco/Golang-sample17-WebAPI-CRUD-for-MongoDb_deploy_to_Docker_Desktop/assets/32194879/9fd5db77-6a25-4fcf-a6ef-988243b922ee)

## 8. Verify your application with Postman

We first create a new book with a **POST** request

http://localhost:8080/books

```json
{
    "title": "The Go Programming Language",
    "author": "Alan A. A. Donovan and Brian W. Kernighan",
    "isbn": "0134190440"
}
```

![image](https://github.com/luiscoco/Golang-sample17-WebAPI-CRUD-for-MongoDb_deploy_to_Docker_Desktop/assets/32194879/c8048b27-e8d1-4dd3-beef-2f51c10d381d)

We verify with a **GET** request the book we created in the previous step

http://localhost:8080/books

![image](https://github.com/luiscoco/Golang-sample17-WebAPI-CRUD-for-MongoDb_deploy_to_Docker_Desktop/assets/32194879/83756fad-adc0-4308-aa57-b246dabe8147)
