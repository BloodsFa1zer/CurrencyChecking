# Use the official Golang image as the base image
FROM golang:1.21-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download and cache dependencies
RUN go mod download

# Copy the entire source code to the working directory
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Use a smaller image for the runtime
FROM alpine:latest

# Install PostgreSQL client and DNS tools
RUN apk --no-cache add postgresql-client bind-tools

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Copy the .env file
COPY .env ./.env

# Expose the port the app runs on
EXPOSE 8080

# Run the application
CMD ["./main"]


## Build
##
FROM golang:1.20-alpine as dev-env

# Copy application data into image
COPY . /Users/mishashevnuk/GolandProjects/app2.4
WORKDIR /Users/mishashevnuk/GolandProjects/app2.4


COPY . .

RUN go mod download

# Copy only .go files, if you want all files to be copied then replace with `COPY . . for the code below.

# Build our application.
# RUN go build -o /go/src/bartmika/mullberry-backend/bin/mullberry-backend
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -o /server

##
## Deploy
##
FROM alpine:latest
RUN mkdir /data

COPY --from=dev-env /server ./
COPY .env ./.env

CMD ["./main"]

