# Build stage
FROM golang:1.20.5-alpine3.17 AS build
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.17
WORKDIR /app
COPY --from=build /app/main .
COPY config.env .

# The exposed port on which the server will run
EXPOSE 8080

# Default command to run when the container starts
CMD ["/app/main"]
