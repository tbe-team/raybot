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

WORKDIR /app

COPY --from=builder /app/bin/raybot /app/raybot

EXPOSE 3000
EXPOSE 60000

CMD ["/app/raybot"]
