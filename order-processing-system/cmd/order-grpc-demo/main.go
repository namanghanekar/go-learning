package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"order-processing-system/internal/grpcx"
	"order-processing-system/pkg/contracts"
)

func main() {
	addr := flag.String("addr", getenv("ORDER_JSON_GRPC_ADDR", "localhost:50061"), "order service legacy JSON gRPC address (host:port)")
	scenario := flag.String("scenario", "success", "success | payment_denied | shipping_refund")
	userID := flag.String("user_id", "", "if empty, a unique user id is generated")
	flag.Parse()

	if *userID == "" {
		*userID = fmt.Sprintf("user_grpc_%d", time.Now().UnixNano())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	conn, err := grpcx.Dial(ctx, *addr)
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	client := contracts.NewOrderServiceClient(conn)

	cartKey := fmt.Sprintf("cart-%s", *userID)
	_, err = client.CreateOrder(ctx, &contracts.CreateOrderRequest{
		IdempotencyKey: cartKey,
		UserID:         *userID,
		Items:          []contracts.OrderItem{{SKU: "SKU-BOOK", Quantity: 2}},
		AmountCents:    2599,
		Email:          "grpc-demo@example.com",
	})
	if err != nil {
		log.Fatalf("CreateOrder: %v", err)
	}

	got, err := client.GetCart(ctx, &contracts.GetCartRequest{UserID: *userID})
	if err != nil {
		log.Fatalf("GetCart: %v", err)
	}
	fmt.Printf("GetCart: order_id=%s status=%s amount_cents=%d items=%d\n", got.OrderID, got.Status, got.AmountCents, len(got.Items))

	checkKey := fmt.Sprintf("checkout-%s", *userID)
	req := &contracts.CheckoutRequest{
		IdempotencyKey: checkKey,
		UserID:         *userID,
		PaymentToken:   "tok_success",
		RecipientName:  "gRPC Demo",
		Address:        "221B Baker Street",
		PaymentOutcome: "success",
	}

	switch *scenario {
	case "success":
	case "payment_denied":
		req.PaymentOutcome = "denied"
	case "shipping_refund":
		req.Address = ""
		req.PaymentOutcome = "success"
	default:
		log.Fatalf("unknown scenario %q (use success, payment_denied, shipping_refund)", *scenario)
	}

	out, err := client.Checkout(ctx, req)
	if err != nil {
		log.Fatalf("Checkout: %v", err)
	}
	fmt.Printf("Checkout: status=%s message=%s payment_id=%s payment_status=%s shipment_id=%s refund_issued=%v\n",
		out.Status, out.Message, out.PaymentID, out.PaymentStatus, out.ShipmentID, out.RefundIssued)
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
