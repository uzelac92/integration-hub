FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o integration-hub ./cmd/integration-hub

FROM alpine:3.19
WORKDIR /app

COPY --from=builder /app/integration-hub /app/integration-hub

COPY config ./config
COPY internal/storage/db/schema ./internal/storage/db/schema

EXPOSE 8080
CMD ["./integration-hub"]