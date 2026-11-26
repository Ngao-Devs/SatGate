FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go-api/go.mod go-api/go.sum ./
RUN go mod download

COPY go-api/main.go ./
COPY go-api/migrations/ ./migrations/

RUN CGO_ENABLED=0 GOOS=linux go build -o /go-api ./main.go

FROM alpine:3.19

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /go-api ./go-api

EXPOSE 8080

CMD ["./go-api"]
