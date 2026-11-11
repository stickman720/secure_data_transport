# Stage 1: Builder
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . . 
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# Stage 2: Runtime
FROM alpine:3.23

WORKDIR /app
COPY --from=builder /app/server ./server

CMD ["./server"]

