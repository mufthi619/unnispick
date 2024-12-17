FROM golang:1.22.7-alpine AS builder

RUN apk add --no-cache gcc musl-dev git

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go

FROM alpine:3.18
RUN apk add --no-cache ca-certificates tzdata netcat-openbsd postgresql-client
RUN adduser -D -g '' appuser
WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/config ./config
COPY scripts/docker-entrypoint.sh .

RUN chmod +x docker-entrypoint.sh && \
    chown -R appuser:appuser /app

USER appuser

ENTRYPOINT ["./docker-entrypoint.sh"]