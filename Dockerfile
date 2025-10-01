# Stage 1: build Go binary
FROM golang:1.25-alpine AS build
WORKDIR /app

RUN apk update && apk upgrade

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
RUN go build -o backend-app ./cmd

# Stage 2: minimal runtime
FROM alpine:3.18
WORKDIR /app
RUN apk update && apk upgrade

COPY --from=build /app/backend-app .


EXPOSE 8080
CMD ["./backend-app"]
