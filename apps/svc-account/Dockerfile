FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY packages/ packages/
COPY apps/svc-account/ apps/svc-account/

WORKDIR /app/apps/svc-account
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/svc-account ./cmd/server

FROM alpine:3.21

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/bin/svc-account .
COPY apps/svc-account/migrations/ ./migrations/

EXPOSE 50052

CMD ["./svc-account"]
