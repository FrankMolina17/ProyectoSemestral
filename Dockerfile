# Etapa 1: builder
FROM golang:1.24-alpine AS builder

WORKDIR /app

ENV GOTOOLCHAIN=auto

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/api ./cmd/api

# Etapa 2: runner
FROM alpine:3.20

COPY --from=builder /bin/api /app/api

ENTRYPOINT ["/app/api"]
