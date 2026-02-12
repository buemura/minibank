FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY packages/ packages/
COPY apps/svc-transaction/ apps/svc-transaction/

WORKDIR /app/apps/svc-transaction
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/svc-transaction ./cmd/server

FROM alpine:3.21

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/bin/svc-transaction .
COPY apps/svc-transaction/migrations/ ./migrations/

EXPOSE 50053

CMD ["./svc-transaction"]
