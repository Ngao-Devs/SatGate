FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go-api ./main.go

FROM alpine:3.19

RUN apk add --no-cache ca-certificates

WORKDIR /root

COPY --from=builder /app/go-api .

EXPOSE 8080

CMD ["/app/go-api"]