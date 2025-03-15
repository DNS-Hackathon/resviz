FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server

# Use minimal image for production
FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/server .

CMD ["./server"]