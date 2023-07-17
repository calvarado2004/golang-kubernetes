# Start from the Go base image
FROM  --platform=linux/amd64 golang:latest as builder

# Add Maintainer Info
LABEL maintainer="Carlos Alvarado carlos-alvarado@outlook.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o main .

# Expose port 8080 to the outside world
RUN chmod +x /app/main

FROM --platform=linux/amd64 busybox:latest

RUN mkdir /app

EXPOSE 8080

COPY --from=builder /app/main /app

CMD [ "/app/main"]