package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"order-processing-system/gen/orderspb"
	"order-processing-system/internal/config"
	"order-processing-system/internal/grpcx"
	"order-processing-system/internal/kafkax"
	"order-processing-system/internal/logx"
	"order-processing-system/internal/postgres"
	"order-processing-system/internal/redisx"
	"order-processing-system/internal/retry"
	"order-processing-system/internal/shutdown"
	"order-processing-system/pkg/contracts"
)

type server struct {
	orderspb.UnimplementedOrderServiceServer
	db        *pgxpool.Pool
	redis     *redis.Client
	producer  *kafkax.Producer
	inventory *contracts.InventoryServiceClient
	payment   *contracts.PaymentServiceClient
	shipping  *contracts.ShippingServiceClient
	log       *slog.Logger
}

type orderPayload struct {
	UserID       string                `json:"user_id"`
	Items        []contracts.OrderItem `json:"items"`
	AmountCents  int64                 `json:"amount_cents"`
	Email        string                `json:"email"`
	FailureCause string                `json:"failure_cause,omitempty"`
}

func (s *server) saveCart(ctx context.Context, req *contracts.CreateOrderRequest) (*contracts.CreateOrderResponse, error) {
	if req.IdempotencyKey == "" {
		req.IdempotencyKey = uuid.NewString()
	}
	if err := validateCartRequest(req); err != nil {
		return nil, err
	}
	if existing, err := s.redis.Get(ctx, "idempotency:cart:"+req.IdempotencyKey).Result(); err == nil && existing != "" {
		return &contracts.CreateOrderResponse{OrderID: existing, UserID: req.UserID, Status: "CART_DUPLICATE_ACCEPTED", Message: "cart idempotent replay"}, nil
	}

	orderID, err := s.findOpenCart(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	if orderID == "" {
		orderID = "ord_" + uuid.NewString()
		if err := s.insertOrder(ctx, orderID, req, "CART"); err != nil {
			return nil, err
		}
	} else if err := s.updateCart(ctx, orderID, req); err != nil {
		return nil, err
	}
	_ = s.redis.Set(ctx, "idempotency:cart:"+req.IdempotencyKey, orderID, 24*time.Hour).Err()
	s.publishAsync(kafkax.TopicOrderEvents, "CartUpdated", orderID, orderPayload{UserID: req.UserID, Items: req.Items, AmountCents: req.AmountCents, Email: req.Email})
	s.log.Info("cart saved", "order_id", orderID, "user_id", req.UserID)
	return &contracts.CreateOrderResponse{OrderID: orderID, UserID: req.UserID, Status: "CART", Message: "cart saved"}, nil
}

func (s *server) loadUserCart(ctx context.Context, req *contracts.GetCartRequest) (*contracts.GetCartResponse, error) {
	if req.UserID == "" {
		return nil, errors.New("user_id is required")
	}
	orderID, err := s.findOpenCart(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	if orderID == "" {
		return &contracts.GetCartResponse{UserID: req.UserID, Status: "EMPTY", Message: "cart is empty"}, nil
	}
	return s.loadCart(ctx, orderID)
}

func (s *server) runCheckout(ctx context.Context, req *contracts.CheckoutRequest) (*contracts.CheckoutResponse, error) {
	if req.UserID == "" {
		return nil, errors.New("user_id is required")
	}
	if req.IdempotencyKey == "" {
		req.IdempotencyKey = uuid.NewString()
	}
	if existing, err := s.redis.Get(ctx, "idempotency:checkout:"+req.IdempotencyKey).Result(); err == nil && existing != "" {
		cart, err := s.loadCart(ctx, existing)
		if err != nil {
			return nil, err
		}
		return &contracts.CheckoutResponse{OrderID: existing, Status: cart.Status, Message: "checkout idempotent replay"}, nil
	}

	orderID, err := s.findOpenCart(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	if orderID == "" {
		return &contracts.CheckoutResponse{Status: "CART_EMPTY", Message: "add items to cart before checkout"}, nil
	}
	cart, err := s.loadCart(ctx, orderID)
	if err != nil {
		return nil, err
	}
	_ = s.redis.Set(ctx, "idempotency:checkout:"+req.IdempotencyKey, orderID, 24*time.Hour).Err()
	if err := s.updateOrder(ctx, orderID, "CHECKOUT_STARTED"); err != nil {
		return nil, err
	}
	s.publishAsync(kafkax.TopicOrderEvents, "CheckoutStarted", orderID, orderPayload{UserID: cart.UserID, Items: cart.Items, AmountCents: cart.AmountCents, Email: cart.Email})

	reservation, err := s.inventory.ReserveInventory(ctx, &contracts.ReserveInventoryRequest{OrderID: orderID, Items: cart.Items})
	if err != nil {
		return s.checkoutFail(ctx, orderID, cart, "INVENTORY_FAILED", err)
	}

	itemCount := totalItemQuantity(cart.Items)
	payment, err := s.payment.ProcessPayment(ctx, &contracts.ProcessPaymentRequest{
		OrderID:        orderID,
		AmountCents:    cart.AmountCents,
		PaymentToken:   req.PaymentToken,
		Address:        req.Address,
		RecipientName:  req.RecipientName,
		ItemCount:      itemCount,
		PaymentOutcome: req.PaymentOutcome,
	})
	if err != nil {
		_, _ = s.inventory.ReleaseInventory(ctx, &contracts.ReleaseInventoryRequest{ReservationID: reservation.ReservationID})
		_ = s.persistPaymentDeclined(ctx, orderID)
		return s.checkoutFail(ctx, orderID, cart, "PAYMENT_FAILED", err)
	}

	shipment, err := s.shipping.CreateShipment(ctx, &contracts.CreateShipmentRequest{OrderID: orderID, Address: req.Address})
	if err != nil {
		refundResp, _ := s.payment.RefundPayment(ctx, &contracts.RefundPaymentRequest{PaymentID: payment.PaymentID})
		_, _ = s.inventory.ReleaseInventory(ctx, &contracts.ReleaseInventoryRequest{ReservationID: reservation.ReservationID})
		_ = s.persistAfterRefund(ctx, orderID, payment.PaymentID)
		refunded := refundResp != nil && refundResp.Refunded
		message := cleanErrorMessage(err)
		if refunded {
			message += "; payment refunded and inventory released"
		}
		return s.checkoutFailWithRefund(ctx, orderID, cart, "SHIPPING_FAILED", errors.New(message), refunded)
	}

	if err := s.persistOrderCompleted(ctx, orderID, payment.PaymentID, shipment.ShipmentID, req.RecipientName, req.Address, itemCount); err != nil {
		return nil, err
	}
	payload := orderPayload{UserID: cart.UserID, Items: cart.Items, AmountCents: cart.AmountCents, Email: cart.Email}
	s.publishAsync(kafkax.TopicOrderEvents, "OrderCompleted", orderID, payload)
	s.publishAsync(kafkax.TopicNotificationRequests, "NotificationRequested", orderID, payload)
	s.log.Info("order completed", "order_id", orderID, "reservation_id", reservation.ReservationID, "payment_id", payment.PaymentID, "shipment_id", shipment.ShipmentID)
	return &contracts.CheckoutResponse{
		OrderID:       orderID,
		Status:        "COMPLETED",
		Message:       "order completed",
		PaymentID:     payment.PaymentID,
		PaymentStatus: payment.PaymentStatus,
		ShipmentID:    shipment.ShipmentID,
	}, nil
}

func (s *server) checkoutFailWithRefund(ctx context.Context, orderID string, cart *contracts.GetCartResponse, status string, cause error, refunded bool) (*contracts.CheckoutResponse, error) {
	_ = s.updateOrder(ctx, orderID, status)
	message := cleanErrorMessage(cause)
	s.publishAsync(kafkax.TopicOrderEvents, "OrderFailed", orderID, orderPayload{UserID: cart.UserID, Items: cart.Items, AmountCents: cart.AmountCents, Email: cart.Email, FailureCause: message})
	s.log.Error("order failed", "order_id", orderID, "status", status, "error", cause)
	return &contracts.CheckoutResponse{OrderID: orderID, Status: status, Message: message, RefundIssued: refunded}, nil
}

func (s *server) checkoutFail(ctx context.Context, orderID string, cart *contracts.GetCartResponse, status string, cause error) (*contracts.CheckoutResponse, error) {
	dbStatus := status
	if status == "INVENTORY_FAILED" || status == "PAYMENT_FAILED" {
		dbStatus = "CART"
	}
	_ = s.updateOrder(ctx, orderID, dbStatus)
	message := cleanErrorMessage(cause)
	s.publishAsync(kafkax.TopicOrderEvents, "OrderFailed", orderID, orderPayload{UserID: cart.UserID, Items: cart.Items, AmountCents: cart.AmountCents, Email: cart.Email, FailureCause: message})
	s.log.Error("order failed", "order_id", orderID, "status", status, "error", cause)
	return &contracts.CheckoutResponse{OrderID: orderID, Status: status, Message: message}, nil
}

func (s *server) CreateOrder(ctx context.Context, req *orderspb.CreateOrderRequest) (*orderspb.CreateOrderResponse, error) {
	cReq := &contracts.CreateOrderRequest{
		IdempotencyKey: req.GetIdempotencyKey(),
		UserID:         req.GetUserId(),
		Items:          pbItemsToContracts(req.GetItems()),
		AmountCents:    req.GetAmountCents(),
		Email:          req.GetEmail(),
	}
	resp, err := s.saveCart(ctx, cReq)
	if err != nil {
		return nil, err
	}
	return &orderspb.CreateOrderResponse{
		OrderId: resp.OrderID,
		UserId:  resp.UserID,
		Status:  resp.Status,
		Message: resp.Message,
	}, nil
}

func (s *server) GetCart(ctx context.Context, req *orderspb.GetCartRequest) (*orderspb.GetCartResponse, error) {
	resp, err := s.loadUserCart(ctx, &contracts.GetCartRequest{UserID: req.GetUserId()})
	if err != nil {
		return nil, err
	}
	return &orderspb.GetCartResponse{
		OrderId:     resp.OrderID,
		UserId:      resp.UserID,
		Status:      resp.Status,
		Items:       contractsItemsToPB(resp.Items),
		AmountCents: resp.AmountCents,
		Email:       resp.Email,
		Message:     resp.Message,
	}, nil
}

func (s *server) Checkout(ctx context.Context, req *orderspb.CheckoutRequest) (*orderspb.CheckoutResponse, error) {
	cReq := &contracts.CheckoutRequest{
		IdempotencyKey: req.GetIdempotencyKey(),
		UserID:         req.GetUserId(),
		PaymentToken:   req.GetPaymentToken(),
		Address:        req.GetAddress(),
		RecipientName:  req.GetRecipientName(),
		PaymentOutcome: req.GetPaymentOutcome(),
	}
	resp, err := s.runCheckout(ctx, cReq)
	if err != nil {
		return nil, err
	}
	return &orderspb.CheckoutResponse{
		OrderId:       resp.OrderID,
		Status:        resp.Status,
		Message:       resp.Message,
		PaymentId:     resp.PaymentID,
		PaymentStatus: resp.PaymentStatus,
		ShipmentId:    resp.ShipmentID,
		RefundIssued:  resp.RefundIssued,
	}, nil
}

func pbItemsToContracts(items []*orderspb.OrderItem) []contracts.OrderItem {
	out := make([]contracts.OrderItem, 0, len(items))
	for _, it := range items {
		if it == nil {
			continue
		}
		out = append(out, contracts.OrderItem{SKU: it.GetSku(), Quantity: it.GetQuantity()})
	}
	return out
}

func contractsItemsToPB(items []contracts.OrderItem) []*orderspb.OrderItem {
	out := make([]*orderspb.OrderItem, 0, len(items))
	for _, it := range items {
		out = append(out, &orderspb.OrderItem{Sku: it.SKU, Quantity: it.Quantity})
	}
	return out
}

type jsonShim struct{ s *server }

func (j jsonShim) CreateOrder(ctx context.Context, req *contracts.CreateOrderRequest) (*contracts.CreateOrderResponse, error) {
	return j.s.saveCart(ctx, req)
}
func (j jsonShim) GetCart(ctx context.Context, req *contracts.GetCartRequest) (*contracts.GetCartResponse, error) {
	return j.s.loadUserCart(ctx, req)
}
func (j jsonShim) Checkout(ctx context.Context, req *contracts.CheckoutRequest) (*contracts.CheckoutResponse, error) {
	return j.s.runCheckout(ctx, req)
}

func cleanErrorMessage(err error) string {
	if err == nil {
		return ""
	}
	if s, ok := status.FromError(err); ok {
		return s.Message()
	}
	return err.Error()
}

func (s *server) insertOrder(ctx context.Context, orderID string, req *contracts.CreateOrderRequest, status string) error {
	items, _ := json.Marshal(req.Items)
	_, err := s.db.Exec(ctx, `insert into orders(id, user_id, status, amount_cents, items, email, created_at, updated_at) values($1,$2,$3,$4,$5,$6,now(),now())`,
		orderID, req.UserID, status, req.AmountCents, items, req.Email)
	return err
}

func (s *server) updateOrder(ctx context.Context, orderID string, status string) error {
	_, err := s.db.Exec(ctx, `update orders set status=$2, updated_at=now() where id=$1`, orderID, status)
	return err
}

func (s *server) updateCart(ctx context.Context, orderID string, req *contracts.CreateOrderRequest) error {
	items, _ := json.Marshal(req.Items)
	_, err := s.db.Exec(ctx, `update orders set items=$2, amount_cents=$3, email=$4, status='CART', updated_at=now() where id=$1`,
		orderID, items, req.AmountCents, req.Email)
	return err
}

func (s *server) findOpenCart(ctx context.Context, userID string) (string, error) {
	var orderID string
	err := s.db.QueryRow(ctx, `select id from orders where user_id=$1 and status='CART' order by updated_at desc limit 1`, userID).Scan(&orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return orderID, nil
}

func (s *server) loadCart(ctx context.Context, orderID string) (*contracts.GetCartResponse, error) {
	var resp contracts.GetCartResponse
	var rawItems []byte
	err := s.db.QueryRow(ctx, `select id, user_id, status, amount_cents, items, email from orders where id=$1`, orderID).
		Scan(&resp.OrderID, &resp.UserID, &resp.Status, &resp.AmountCents, &rawItems, &resp.Email)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(rawItems, &resp.Items); err != nil {
		return nil, err
	}
	resp.Message = "cart loaded"
	return &resp, nil
}

func validateCartRequest(req *contracts.CreateOrderRequest) error {
	if req.UserID == "" {
		return errors.New("user_id is required")
	}
	if len(req.Items) == 0 {
		return errors.New("at least one item is required")
	}
	if req.AmountCents <= 0 {
		return errors.New("amount_cents must be positive")
	}
	for _, item := range req.Items {
		if item.SKU == "" || item.Quantity <= 0 {
			return errors.New("each item needs sku and positive quantity")
		}
	}
	return nil
}

func (s *server) publish(ctx context.Context, topic string, eventType string, orderID string, payload any) error {
	data, _ := json.Marshal(payload)
	return s.producer.Publish(ctx, topic, orderID, kafkax.Event{ID: uuid.NewString(), Type: eventType, OrderID: orderID, Payload: data})
}

func (s *server) publishAsync(topic string, eventType string, orderID string, payload any) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		if err := s.publish(ctx, topic, eventType, orderID, payload); err != nil {
			s.log.Error("event publish failed", "topic", topic, "event_type", eventType, "order_id", orderID, "error", err)
		}
	}()
}

func totalItemQuantity(items []contracts.OrderItem) int32 {
	var n int32
	for _, it := range items {
		n += it.Quantity
	}
	return n
}

func (s *server) persistPaymentDeclined(ctx context.Context, orderID string) error {
	_, err := s.db.Exec(ctx, `update orders set payment_status='DECLINED', updated_at=now() where id=$1`, orderID)
	return err
}

func (s *server) persistAfterRefund(ctx context.Context, orderID, paymentID string) error {
	_, err := s.db.Exec(ctx, `update orders set payment_id=$2, payment_status='REFUNDED', shipment_id=null, updated_at=now() where id=$1`, orderID, paymentID)
	return err
}

func (s *server) persistOrderCompleted(ctx context.Context, orderID, paymentID, shipmentID, recipient, address string, itemCount int32) error {
	_, err := s.db.Exec(ctx, `update orders set status='COMPLETED', payment_id=$2, payment_status='SUCCESS', shipment_id=$3, recipient_name=$4, shipping_address=$5, item_count=$6, updated_at=now() where id=$1`,
		orderID, paymentID, shipmentID, recipient, address, itemCount)
	return err
}

func main() {
	cfg := config.Load("order-service", ":50051")
	log := logx.New(cfg.ServiceName)
	ctx := shutdown.Context()
	db, err := postgres.Connect(ctx, cfg.PostgresURL)
	if err != nil {
		log.Error("postgres connection failed", "error", err)
		return
	}
	rdb, err := redisx.Connect(ctx, cfg.RedisAddr)
	if err != nil {
		log.Error("redis connection failed", "error", err)
		return
	}
	inventoryConn, err := dialWithRetry(ctx, env("INVENTORY_GRPC_ADDR", "localhost:50052"))
	if err != nil {
		log.Error("inventory dial failed", "error", err)
		return
	}
	paymentConn, err := dialWithRetry(ctx, env("PAYMENT_GRPC_ADDR", "localhost:50053"))
	if err != nil {
		log.Error("payment dial failed", "error", err)
		return
	}
	shippingConn, err := dialWithRetry(ctx, env("SHIPPING_GRPC_ADDR", "localhost:50054"))
	if err != nil {
		log.Error("shipping dial failed", "error", err)
		return
	}
	defer inventoryConn.Close()
	defer paymentConn.Close()
	defer shippingConn.Close()

	listener, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		log.Error("listen failed", "error", err)
		return
	}
	producer := kafkax.NewProducer(cfg.KafkaBrokers, log)
	defer producer.Close()

	srv := &server{
		db: db, redis: rdb, producer: producer,
		inventory: contracts.NewInventoryServiceClient(inventoryConn),
		payment:   contracts.NewPaymentServiceClient(paymentConn),
		shipping:  contracts.NewShippingServiceClient(shippingConn),
		log:       log,
	}

	grpcServer := grpc.NewServer()
	orderspb.RegisterOrderServiceServer(grpcServer, srv)
	reflection.Register(grpcServer)

	legacyAddr := env("ORDER_JSON_GRPC_ADDR", ":50061")
	legacyListener, legacyErr := net.Listen("tcp", legacyAddr)
	var legacyServer *grpc.Server
	if legacyErr == nil {
		legacyServer = grpcx.Server()
		contracts.RegisterOrderServiceServer(legacyServer, jsonShim{s: srv})
		reflection.Register(legacyServer)
		go func() {
			log.Info("legacy json-grpc listening", "addr", legacyAddr)
			_ = legacyServer.Serve(legacyListener)
		}()
	} else {
		log.Warn("legacy json-grpc disabled (listen failed)", "addr", legacyAddr, "error", legacyErr)
	}

	go func() {
		<-ctx.Done()
		if legacyServer != nil {
			legacyServer.GracefulStop()
		}
		grpcServer.GracefulStop()
	}()
	log.Info("order service listening", "addr", cfg.GRPCAddr)
	if err := grpcServer.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		log.Error("grpc server failed", "error", err)
	}
}

func dialWithRetry(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	var conn *grpc.ClientConn
	err := retry.Do(ctx, 20, 250*time.Millisecond, func(ctx context.Context) error {
		c, err := grpcx.Dial(ctx, addr)
		if err != nil {
			return err
		}
		conn = c
		return nil
	})
	return conn, err
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
