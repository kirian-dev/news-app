FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o news-app ./cmd/server

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/news-app .

COPY templates ./templates

CMD ["./news-app"]
