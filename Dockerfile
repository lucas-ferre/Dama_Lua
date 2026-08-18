FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/damas ./cmd/damas

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/damas /app/damas

ENTRYPOINT ["/app/damas"]
