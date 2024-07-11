# Start from the golang base image
FROM golang:1.22 AS builder

# Add Maintainer Info
LABEL maintainer="Russell Jones <jonesrussell42@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Define a variable for the binary name
ARG BINARY_NAME=gojekyll

# Build the Go app
RUN CGO_ENABLED=0 go build -o ${BINARY_NAME} .

# Start a new stage from scratch
FROM alpine:latest  

RUN apk --no-cache add ca-certificates ncurses

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/${BINARY_NAME} .
# Copy the static directory from the previous stage
COPY --from=builder /app/static ./static

RUN chmod u+x /root/${BINARY_NAME}

# Expose port 3000 to the outside world
EXPOSE 3000

# Command to run the executable
CMD ["/root/${BINARY_NAME}", "website"]
