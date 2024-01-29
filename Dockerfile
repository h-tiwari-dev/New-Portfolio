# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project to the working directory
COPY . .

# Build the application
RUN go build -o main .

# Expose the port your application will run on
EXPOSE 8080

# Command to run your application with MySQL credentials from .env file
CMD ["./main"]

