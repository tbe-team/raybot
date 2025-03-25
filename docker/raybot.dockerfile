FROM node:20 AS ui-builder

WORKDIR /app/ui

COPY ui/package.json ui/pnpm-lock.yaml ./
RUN corepack enable && pnpm install --frozen-lockfile

COPY ui ./
RUN pnpm run build


FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=ui-builder /app/ui/dist ui/dist

RUN CGO_ENABLED=1 go build -o bin/raybot cmd/raybot/main.go


FROM debian:bookworm-slim AS prod


RUN mkdir -p /app && \
    groupadd -r -g 1000 raybot && \
    useradd -r -m -u 1000 -g raybot raybot && \
    chown -R raybot:raybot /app

WORKDIR /app

COPY --from=builder /app/bin/raybot /app/raybot

USER raybot

# raybot will occupy these ports:
#  - 3000: http port
#  - 60000: grpc port
EXPOSE 3000 60000

CMD ["/app/raybot"]
