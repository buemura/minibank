FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY packages/ packages/
COPY apps/svc-auth/ apps/svc-auth/

WORKDIR /app/apps/svc-auth
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/svc-auth ./cmd/server

FROM alpine:3.21

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/bin/svc-auth .
COPY apps/svc-auth/migrations/ ./migrations/

EXPOSE 50051

CMD ["./svc-auth"]
