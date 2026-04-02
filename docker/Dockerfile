# ── Builder stage ─────────────────────────────────────────────────────────────
FROM golang:1.22-bookworm AS builder

WORKDIR /build
COPY go.mod go.sum* ./
RUN go mod download 2>/dev/null || true

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /build/vector-editor ./cmd/server

# ── Runtime stage ─────────────────────────────────────────────────────────────
FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        inkscape \
        fontconfig \
        fonts-noto \
        fonts-roboto \
        fonts-open-sans \
        ca-certificates && \
    fc-cache -fv && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /build/vector-editor /app/vector-editor

RUN mkdir -p /app/data

ENV VEC_PORT=8092 \
    VEC_DATA_DIR=/app/data \
    VEC_INKSCAPE_PATH=inkscape \
    VEC_MAX_WORKERS=4

EXPOSE 8092

ENTRYPOINT ["/app/vector-editor"]
