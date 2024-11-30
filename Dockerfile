FROM golang:1.23.2 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/rate-limiter/main.go

FROM scratch
COPY --from=builder /app/main /app/
WORKDIR /app
CMD ["./main"]