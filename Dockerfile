# Build stage
FROM golang:1.20.5-alpine3.17 AS build
WORKDIR /app
COPY . .
# Installing required apps for the server to functional properly
# Curl isn't installed by default in the alpine
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
RUN go build -o main main.go

# Run stage
FROM alpine:3.17
WORKDIR /app
COPY --from=build /app/main .
# Copying golang-migrate and others
COPY --from=build /app/migrate .
COPY db/migrations/ ./migrations
COPY wait-for.sh .
COPY startup.sh .
COPY config.env .

# The exposed port on which the server will run
EXPOSE 8080

# Default command to run when the container starts
# This is equal to : ENTRYPOINT CMD --> /app/startup.sh /app/main
CMD ["/app/main"]
ENTRYPOINT ["/app/startup.sh"]
