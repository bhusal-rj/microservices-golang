#base go image
FROM golang:1.23-alpine as builder

RUN mkdir /app

COPY go.mod /app

WORKDIR /app

RUN go mod download

COPY . /app

RUN CGO_ENABLED=0 go build -o frontendApp ./cmd/web

RUN chmod +x /app/frontendApp

FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/frontendApp /app

CMD ["/app/frontendApp"]
