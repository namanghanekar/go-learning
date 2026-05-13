FROM golang:1.24

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o notification-service ./cmd/main.go

EXPOSE 8084

CMD ["./notification-service"]