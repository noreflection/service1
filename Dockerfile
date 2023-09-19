# Use the official Golang image as the base image
FROM golang:1.16

# Set the working directory inside the container
WORKDIR /app

# Copy the Go application code into the container
COPY main.go .

# Build the Go application
RUN go build -o main .

# Expose the port your Go application listens on (if applicable)
EXPOSE 8080

# Command to run the Go application
CMD ["./main"]