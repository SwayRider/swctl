# syntax=docker/dockerfile:1.4
# Build context must be the repo root (go.work covers all modules).
FROM --platform=$BUILDPLATFORM golang:latest AS builder

ARG TARGETOS
ARG TARGETARCH

ENV CGO_ENABLED=0
ENV GOOS=${TARGETOS}
ENV GOARCH=${TARGETARCH}

WORKDIR /build

COPY go.work go.work.sum ./
COPY swlib/ swlib/
COPY grpcclients/ grpcclients/
COPY protos/ protos/
COPY swctl/ swctl/

RUN go build -o /swctl ./swctl/cmd/swctl

FROM --platform=$TARGETPLATFORM debian:bookworm-slim
COPY --from=builder /swctl /swctl
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["/swctl"]
