# syntax=docker/dockerfile:1.4
FROM --platform=$BUILDPLATFORM golang:latest AS builder

ARG TARGETOS
ARG TARGETARCH

ENV CGO_ENABLED=0
ENV GOOS=${TARGETOS}
ENV GOARCH=${TARGETARCH}

WORKDIR /app

COPY . .

RUN go clean -modcache && \
    go mod download && \
    go build -o swctl ./cmd/swctl

FROM --platform=$TARGETPLATFORM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/swctl .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["/app/swctl"]
