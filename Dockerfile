# Use the official Golang image to create a build stage.
FROM golang:1.22 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Copy the .env file from the host into the container
COPY .env ./.env

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Set the working directory to the location of your main.go file
WORKDIR /app/cmd/ethereum-tracker-app

# Build the Go app
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o ethereum-tracker-app .

# Start a new stage from scratch
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/cmd/ethereum-tracker-app/ethereum-tracker-app .
COPY --from=builder /app/docs /docs  
COPY --from=builder /app/.env .

RUN chmod +x ethereum-tracker-app

# Expose port 8080 to the outside world
EXPOSE 8000

# Command to run the executable with default arguments
CMD ["./ethereum-tracker-app"]
