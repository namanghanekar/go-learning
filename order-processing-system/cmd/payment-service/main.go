package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net"
	"strings"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"

	"order-processing-system/internal/config"
	"order-processing-system/internal/grpcx"
	"order-processing-system/internal/logx"
	"order-processing-system/internal/redisx"
	"order-processing-system/internal/shutdown"
	"order-processing-system/pkg/contracts"
)

type server struct {
	contracts.PaymentServiceServer
	redis *redis.Client
	log   *slog.Logger
}

func outcomeDenied(outcome string) bool {
	return strings.EqualFold(strings.TrimSpace(outcome), "denied")
}

func outcomeSuccess(outcome string) bool {
	return strings.EqualFold(strings.TrimSpace(outcome), "success")
}

func (s *server) ProcessPayment(ctx context.Context, req *contracts.ProcessPaymentRequest) (*contracts.ProcessPaymentResponse, error) {
	if req.AmountCents <= 0 {
		return nil, errors.New("amount must be positive")
	}
	if req.ItemCount < 0 {
		return nil, errors.New("item_count must be zero or positive")
	}
	if outcomeDenied(req.PaymentOutcome) {
		return nil, errors.New("payment declined (payment_outcome=denied)")
	}
	token := strings.TrimSpace(req.PaymentToken)
	if !outcomeSuccess(req.PaymentOutcome) {
		if token == "" || strings.EqualFold(token, "fail") {
			return nil, errors.New("payment declined")
		}
	}
	key := "payment:order:" + req.OrderID
	if existing, err := s.redis.Get(ctx, key).Result(); err == nil && existing != "" {
		return &contracts.ProcessPaymentResponse{PaymentID: existing, PaymentStatus: "SUCCESS"}, nil
	}
	paymentID := "pay_" + uuid.NewString()
	meta, _ := json.Marshal(map[string]any{
		"order_id":        req.OrderID,
		"amount_cents":    req.AmountCents,
		"address":         req.Address,
		"recipient_name":  req.RecipientName,
		"item_count":      req.ItemCount,
		"payment_outcome": req.PaymentOutcome,
	})
	_ = s.redis.Set(ctx, "payment:meta:"+paymentID, string(meta), 0).Err()
	if err := s.redis.Set(ctx, key, paymentID, 0).Err(); err != nil {
		return nil, err
	}
	s.log.Info("payment processed", "order_id", req.OrderID, "payment_id", paymentID, "amount_cents", req.AmountCents, "item_count", req.ItemCount, "recipient", req.RecipientName)
	return &contracts.ProcessPaymentResponse{PaymentID: paymentID, PaymentStatus: "SUCCESS"}, nil
}

func (s *server) RefundPayment(ctx context.Context, req *contracts.RefundPaymentRequest) (*contracts.RefundPaymentResponse, error) {
	if req.PaymentID == "" {
		return &contracts.RefundPaymentResponse{Refunded: false}, nil
	}
	_ = s.redis.Set(ctx, "refund:"+req.PaymentID, "completed", 0).Err()
	s.log.Info("payment refunded", "payment_id", req.PaymentID)
	return &contracts.RefundPaymentResponse{Refunded: true}, nil
}

func main() {
	cfg := config.Load("payment-service", ":50053")
	log := logx.New(cfg.ServiceName)
	ctx := shutdown.Context()
	rdb, err := redisx.Connect(ctx, cfg.RedisAddr)
	if err != nil {
		log.Error("redis connection failed", "error", err)
		return
	}
	listener, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		log.Error("listen failed", "error", err)
		return
	}
	grpcServer := grpcx.Server()
	contracts.RegisterPaymentServiceServer(grpcServer, &server{redis: rdb, log: log})
	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()
	log.Info("payment service listening", "addr", cfg.GRPCAddr)
	if err := grpcServer.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		log.Error("grpc server failed", "error", err)
	}
}
