# gRPC Testing

## Inventory Service Test

Port:
50051

### Request

Product ID: 101
Quantity: 2

### Expected Response

{
  "available": true
}

---

## Payment Service Test

Port:
50052

### Request

Order ID: 1
Amount: 500

### Expected Response

{
  "success": true
}