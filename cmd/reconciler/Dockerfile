FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o reconciler ./cmd/reconciler

FROM alpine:3.19
WORKDIR /app

COPY --from=builder /app/reconciler .

COPY --from=builder /app/config ./config

CMD ["./reconciler"]