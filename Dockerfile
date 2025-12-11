# ============================
# 1) Build stage
# ============================
FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY . .

# Статическая сборка (лучше для Docker)
ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go build -o server .

# ============================
# 2) Final stage
# ============================
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8081

# os.Args[1] = callback URL
# os.Args[2] = AI URL
ENTRYPOINT ["./server"]
