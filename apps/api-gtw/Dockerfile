FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY packages/ packages/
COPY apps/api-gtw/ apps/api-gtw/

WORKDIR /app/apps/api-gtw
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api-gtw ./cmd/server

FROM alpine:3.21

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/bin/api-gtw .

EXPOSE 8080

CMD ["./api-gtw"]
