FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o pr-reviewer ./cmd/app

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/pr-reviewer /app/pr-reviewer

EXPOSE 8080

ENV DB_DSN=postgres://postgres:postgres@db:5432/postgres?sslmode=disable

CMD ["/app/pr-reviewer"]
