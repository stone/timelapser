# Build stage
FROM golang:alpine AS builder

RUN apk add --no-cache git make
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o timelapser .

# Final stage
FROM alpine:3

LABEL org.opencontainers.image.source="https://github.com/stone/timelapser" \
      org.opencontainers.image.description="Create timelapses from http based camera snapshots" \
      org.opencontainers.image.licenses="MIT" \
      org.opencontainers.image.version="1.0.0"

# Install ffmpeg and required runtime dependencies
RUN apk add --no-cache ffmpeg ffmpeg-libs ca-certificates tzdata && \
  adduser -D -H -h /app timelapser
WORKDIR /app
COPY --from=builder /app/timelapser .
COPY --from=builder /app/example-config.yaml .
RUN chown -R timelapser:timelapser /app  && \
    chmod +x /app/timelapser
USER timelapser
ENTRYPOINT ["/app/timelapser"]
CMD ["-config", "/app/example-config.yaml"]
