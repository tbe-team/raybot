FROM node:22 AS ui-builder

WORKDIR /app/ui

COPY ui/package.json ui/pnpm-lock.yaml ./
RUN corepack enable && pnpm install --frozen-lockfile

COPY ui ./
RUN pnpm run build


FROM golang:1.24-bullseye AS builder

ARG PKG_PREFIX
ARG VERSION
ARG BUILD_DATE

RUN apt-get update && \
    apt-get install -y gcc-aarch64-linux-gnu

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=ui-builder /app/ui/dist ui/dist

ENV CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=arm64 \
    CC=aarch64-linux-gnu-gcc

RUN go build -o bin/raybot\
    -ldflags "\
        -X $PKG_PREFIX.Version=$VERSION \
        -X $PKG_PREFIX.Date=$BUILD_DATE" \
    cmd/raybot/main.go


FROM alpine:3.21.3 AS prod

WORKDIR /app

COPY --from=builder /app/bin/raybot /app/raybot
