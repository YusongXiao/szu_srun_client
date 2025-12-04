# Build stage
FROM golang:1.20-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o /srunClient .

# Final image
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /srunClient /srunClient
ENV TZ=UTC
# Do NOT set SRUN_USERNAME/SRUN_PASSWORD here; pass them at runtime for security.
# To auto login when the container starts, set SRUN_AUTO_LOGIN=1.

EXPOSE 0
ENTRYPOINT ["/srunClient"]
