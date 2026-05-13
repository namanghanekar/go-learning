FROM golang:1.24

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o shipment-service ./cmd/main.go

EXPOSE 8083

CMD ["./shipment-service"]