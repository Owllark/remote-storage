# Use the official Golang base image
FROM golang:1.20-alpine

# Set the working directory inside the container
WORKDIR /app

RUN mkdir storage

# Copy the source code into the container
COPY src .

# Build the Go application
RUN go build -o storage-server

# Expose a port if your application listens on a specific port
EXPOSE 8080

# Set the command to run your application by default
CMD ["--memory=2g", "./storage-server"]


