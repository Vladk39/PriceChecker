# Stage 1: build
FROM golang:1.23.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY envs /envs

COPY migrations /migrations

RUN go build -o PriceChecker ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/PriceChecker .

COPY --from=builder /envs /envs
COPY --from=builder /migrations /migrations

EXPOSE 8080

CMD ["./PriceChecker"]
