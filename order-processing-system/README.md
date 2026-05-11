# Distributed Order Processing System

Production-style Go microservices demo covering order creation, inventory reservation, payment processing, shipment creation, and user notification.

## Architecture

- API Gateway exposes HTTP `POST /orders`.
- Order Service owns the order lifecycle and runs a Saga.
- Inventory, Payment, and Shipping are synchronous gRPC services.
- Kafka carries order and notification events.
- Notification Service consumes Kafka with a worker pool.
- PostgreSQL stores orders.
- Redis stores idempotency keys, inventory counters, payment/shipment idempotency, and notification dedupe state.

## Reliability Features

- gRPC inter-service calls.
- Kafka producer and consumer flow.
- Retry with exponential backoff.
- Idempotency for order creation, payments, shipments, and notifications.
- Saga compensation:
  - payment failure releases inventory.
  - shipping failure refunds payment and releases inventory.
- Graceful shutdown through signal-aware contexts.
- Structured JSON logging with `log/slog`.
- Concurrency-safe notification worker pool using goroutines and channels.

## Run Locally

```powershell
cd C:\Users\sreya\OneDrive\Desktop\go-learning\order-processing-system
docker compose up --build
```

Wait until the API gateway logs that it is listening on `:8080`.

## Test Cart And Checkout Flow

First add order details to cart. This does not process payment yet.

```powershell
$cart = @{
  idempotency_key = "cart-demo-1"
  user_id = "user_123"
  items = @(@{ sku = "SKU-BOOK"; quantity = 2 })
  amount_cents = 2599
  email = "sreya@example.com"
} | ConvertTo-Json -Depth 5

Invoke-RestMethod -Method Post -Uri http://localhost:8080/cart -ContentType "application/json" -Body $cart
```

View the cart:

```powershell
Invoke-RestMethod -Method Get -Uri "http://localhost:8080/cart?user_id=user_123"
```

Checkout with payment and address:

```powershell
$checkout = @{
  idempotency_key = "checkout-demo-1"
  user_id = "user_123"
  payment_token = "tok_success"
  address = "221B Baker Street"
} | ConvertTo-Json -Depth 5

Invoke-RestMethod -Method Post -Uri http://localhost:8080/checkout -ContentType "application/json" -Body $checkout
```

Expected response:

```json
{
  "order_id": "ord_...",
  "status": "COMPLETED",
  "message": "order completed"
}
```

`POST /orders` is kept as a compatibility alias for `POST /cart`.

## Old One-Step Order Creation

```powershell
$body = @{
  idempotency_key = "demo-order-1"
  user_id = "user_123"
  items = @(@{ sku = "SKU-BOOK"; quantity = 2 })
  amount_cents = 2599
  payment_token = "tok_success"
  address = "221B Baker Street"
  email = "sreya@example.com"
} | ConvertTo-Json -Depth 5

Invoke-RestMethod -Method Post -Uri http://localhost:8080/orders -ContentType "application/json" -Body $body
```

Expected response:

```json
{
  "order_id": "ord_...",
  "status": "COMPLETED",
  "message": "order completed"
}
```

Run the same request again with the same `idempotency_key`; it should return `DUPLICATE_ACCEPTED`.

## Test Compensation Flow

Payment failure:

```powershell
$body = @{
  idempotency_key = "demo-order-payment-fail"
  user_id = "user_123"
  items = @(@{ sku = "SKU-PHONE"; quantity = 1 })
  amount_cents = 59900
  payment_token = "fail"
  address = "221B Baker Street"
  email = "sreya@example.com"
} | ConvertTo-Json -Depth 5

Invoke-RestMethod -Method Post -Uri http://localhost:8080/orders -ContentType "application/json" -Body $body
```

Expected status: `PAYMENT_FAILED`. Inventory is released by compensation.

Shipping failure:

```powershell
$body = @{
  idempotency_key = "demo-order-shipping-fail"
  user_id = "user_123"
  items = @(@{ sku = "SKU-LAPTOP"; quantity = 1 })
  amount_cents = 99900
  payment_token = "tok_success"
  address = ""
  email = "sreya@example.com"
} | ConvertTo-Json -Depth 5

Invoke-RestMethod -Method Post -Uri http://localhost:8080/orders -ContentType "application/json" -Body $body
```

Expected status: `SHIPPING_FAILED`. Payment refund and inventory release are triggered.

## Developer Commands

```powershell
go test ./...
go build ./...
docker compose logs -f order-service notification-service
docker compose down -v
```

## Ports

- API Gateway: `localhost:8080`
- PostgreSQL: `localhost:5432`
- Redis: `localhost:6379`
- Kafka: internal Docker network plus host `localhost:9092`
