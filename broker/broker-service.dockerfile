#base go image
FROM golang:1.23-alpine as builder

RUN mkdir /app

COPY go.mod /app

WORKDIR /app

RUN go mod download

COPY . /app

RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

RUN chmod +x /app/brokerApp

FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/brokerApp /app

CMD ["/app/brokerApp"]
