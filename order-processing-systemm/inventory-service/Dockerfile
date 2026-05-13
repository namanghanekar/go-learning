FROM golang:1.24

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o inventory-service ./cmd/main.go

EXPOSE 50051

CMD ["./inventory-service"]