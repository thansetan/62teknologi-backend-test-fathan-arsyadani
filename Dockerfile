FROM golang:1.21.4 AS builder

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./62teknologi-be ./cmd/app/

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/62teknologi-be ./62teknologi-be
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/.env ./.env

EXPOSE 8080
CMD "./62teknologi-be"