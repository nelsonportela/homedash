# Minimal Go build for dashboard
FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go mod init homedash && \
    go mod tidy && \
    go build -ldflags="-s -w" -o dashboard main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/dashboard .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
EXPOSE 8080
CMD ["/app/dashboard"]
