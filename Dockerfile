# Stage 1: Base & Caching
FROM golang:1.26-alpine AS base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Stage 2: Development target (Contains Air for live reload)
FROM base AS dev
RUN go install github.com/air-verse/air@latest
COPY . .
CMD ["air", "-c", ".air.toml"]

# Stage 3: Builder (Compiles the static binary)
FROM base AS builder
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/app ./cmd/api

# Stage 4: Production target (Only the tiny binary, no Go SDK)
FROM alpine:3.21 AS prod
COPY --from=builder /bin/app /app
CMD ["/app"]