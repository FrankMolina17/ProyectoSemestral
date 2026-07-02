# Etapa 1: builder
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /bin/api ./cmd/api

# Etapa 2: runner
FROM alpine:3.20

COPY --from=builder /bin/api /app/api

ENTRYPOINT ["/app/api"]
