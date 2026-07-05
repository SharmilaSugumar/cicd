FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /worker ./cmd/worker

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /worker .
# Install Docker CLI so the worker can launch containers, plus git and build tools
RUN apk add --no-cache docker-cli git go nodejs npm python3 py3-pip make g++ gcc rust cargo openjdk17 maven
CMD ["./worker"]
