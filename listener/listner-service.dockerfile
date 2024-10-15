#base go image
FROM golang:1.23-alpine

RUN mkdir /app

COPY go.mod /app

WORKDIR /app

RUN go mod download

COPY . /app


RUN CGO_ENABLED=0 go build -o listnerApp .

RUN chmod +x /app/listnerApp

CMD ["./listnerApp"]
