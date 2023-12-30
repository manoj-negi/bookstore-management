FROM golang:1.21-alpine3.17

# Install necessary packages for Go development
RUN apk add --no-cache git

RUN mkdir /app
ADD . /app
WORKDIR /app

# Copy all files and directories from your project into the Docker image
COPY . /app/

# Build your Go application
RUN go build -o main cmd/api/main.go

# Set the command to run your application
CMD ["/app/cmd/api/main"]