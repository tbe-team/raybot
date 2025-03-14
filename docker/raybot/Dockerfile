FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o bin/raybot cmd/raybot/main.go

FROM golang:1.24 AS prod

WORKDIR /app

COPY --from=builder /app/bin/raybot /app/raybot

EXPOSE 3000
EXPOSE 50051

CMD ["/app/raybot"]
