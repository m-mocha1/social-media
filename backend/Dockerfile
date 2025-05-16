FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o backend .

FROM alpine:latest
WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/backend .

EXPOSE 8080

CMD ["go run ."]
