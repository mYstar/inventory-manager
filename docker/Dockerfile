FROM golang:1.24-alpine AS builder
WORKDIR /build
COPY . .
RUN go build -o main .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /build/main .
EXPOSE 80
CMD ["./main"]