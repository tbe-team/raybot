FROM node:22 AS ui-builder

WORKDIR /app/ui

COPY ui/package.json ui/pnpm-lock.yaml ./
RUN corepack enable && pnpm install --frozen-lockfile

COPY ui ./
RUN pnpm run build


FROM golang:1.24 AS builder

ARG PKG_PREFIX
ARG VERSION
ARG BUILD_DATE

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=ui-builder /app/ui/dist ui/dist

RUN CGO_ENABLED=1 \
	go build -ldflags "\
		-X $PKG_PREFIX.Version=$VERSION \
		-X $PKG_PREFIX.Date=$BUILD_DATE" \
	-o bin/raybot cmd/raybot/main.go


FROM debian:bookworm-slim AS prod


RUN mkdir -p /app && \
    groupadd -r -g 1000 raybot && \
    useradd -r -m -u 1000 -g raybot raybot && \
    chown -R raybot:raybot /app

WORKDIR /app

COPY --from=builder /app/bin/raybot /app/raybot

USER raybot

CMD ["/app/raybot"]
