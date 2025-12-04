# Build stage
FROM golang:alpine AS builder

# Set Go Proxy to Aliyun
ENV GOPROXY=https://mirrors.aliyun.com/goproxy/

WORKDIR /app

# Copy go mod files
COPY go.mod ./
# If you have a go.sum file, uncomment the next line
# COPY go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o srunClient .

# Final stage
FROM alpine:latest

# Set Alpine repositories to Aliyun
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

WORKDIR /app

COPY --from=builder /app/srunClient .

# Run the binary
CMD ["./srunClient"]
