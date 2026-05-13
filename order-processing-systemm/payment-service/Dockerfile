FROM golang:1.24

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o payment-service ./cmd/main.go

EXPOSE 8082

CMD ["./payment-service"]