FROM golang:1.22.7 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o gophertalk ./cmd/gophertalk/main.go

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /root/

COPY --from=builder /app/gophertalk .
COPY --from=builder /app/internal/migrations /migrations

EXPOSE 3000

CMD ["./gophertalk"]
