FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o shortener main.go

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/shortener .
EXPOSE 8080
CMD ["./shortener"]