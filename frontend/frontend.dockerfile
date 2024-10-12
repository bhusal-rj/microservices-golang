#base go image
FROM golang:1.23-alpine as builder

RUN mkdir /app

COPY go.mod /app

WORKDIR /app

RUN go mod download

COPY . /app

RUN ls 

RUN CGO_ENABLED=0 go build -o frontendApp ./cmd/web

RUN chmod +x /app/frontendApp

FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/frontendApp /app
RUN mkdir -p /app/cmd/web/templates
COPY --from=builder /app/cmd/web/templates /app/cmd/web/templates

CMD ["/app/frontendApp"]
