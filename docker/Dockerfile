FROM golang:1.26-alpine AS builder
WORKDIR /build
COPY . .
RUN go build -o main ./cmd/inventory_manager

FROM alpine:3.23.3
WORKDIR /app
COPY --from=builder /build/main .
EXPOSE 8080
CMD ["./main"]