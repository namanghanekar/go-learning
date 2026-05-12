# Order processing

gRPC services: cart -> checkout (inventory, payment, shipping) -> Kafka -> notification worker.

**Ports:** Postman/protobuf `50051`, JSON demo client `50061`, Postgres `5432`, Redis `6379`, Kafka `9092`.

## Run

```powershell
cd C:\Users\sreya\OneDrive\Desktop\go-orderpro\order-processing-system
docker compose up --build
```

Postman: server `localhost:50051`, import `proto/order_processing.proto`, message JSON **snake_case** (`user_id`, `payment_token`, …).

```powershell
go run ./cmd/order-grpc-demo -addr localhost:50061 -scenario success
```

```powershell
docker compose exec postgres psql -U orders -d orders -c "select id,user_id,status,updated_at from orders order by updated_at desc limit 5;"
```

```powershell
go test ./...
docker compose down -v
```

Regenerate protobuf Go code after editing `proto/order_processing.proto`:

```powershell
protoc -I proto --go_out=. --go_opt=module=order-processing-system --go-grpc_out=. --go-grpc_opt=module=order-processing-system proto/order_processing.proto
```
