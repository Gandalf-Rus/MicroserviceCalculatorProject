# Start from golang:1.21-alpine base image
FROM golang:1.21-alpine

# The latest alpine images don't have some tools like (`git` and `bash`).
# Adding git, bash and openssh to the image
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the entire project into the container
COPY . .

# Build the Go app
RUN go build -o orchestratorExe ./orchestrator/cmd/main.go

# Install goose CLI tool
RUN go get -u github.com/pressly/goose/cmd/goose@latest

# Sleep for 30 seconds to wait for PostgreSQL to start
CMD ["sh", "-c", "sleep 30 && goose -dir /app/migrations postgres 'postgres://admin:admin@postgres:5432/MicroserviceCalculatorDB' up && ./orchestratorExe"]

# Expose port 8080 to the outside world
EXPOSE 8080
