package main

import (
	"context"
	"errors"
	"log/slog"
	"net"

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
	contracts.ShippingServiceServer
	redis *redis.Client
	log   *slog.Logger
}

func (s *server) CreateShipment(ctx context.Context, req *contracts.CreateShipmentRequest) (*contracts.CreateShipmentResponse, error) {
	if req.Address == "" {
		return nil, errors.New("shipping address is required")
	}
	key := "shipment:order:" + req.OrderID
	if existing, err := s.redis.Get(ctx, key).Result(); err == nil && existing != "" {
		return &contracts.CreateShipmentResponse{ShipmentID: existing}, nil
	}
	shipmentID := "ship_" + uuid.NewString()
	if err := s.redis.Set(ctx, key, shipmentID, 0).Err(); err != nil {
		return nil, err
	}
	s.log.Info("shipment created", "order_id", req.OrderID, "shipment_id", shipmentID)
	return &contracts.CreateShipmentResponse{ShipmentID: shipmentID}, nil
}

func (s *server) CancelShipment(ctx context.Context, req *contracts.CancelShipmentRequest) (*contracts.CancelShipmentResponse, error) {
	if req.ShipmentID == "" {
		return &contracts.CancelShipmentResponse{Cancelled: false}, nil
	}
	_ = s.redis.Set(ctx, "shipment_cancel:"+req.ShipmentID, "completed", 0).Err()
	s.log.Info("shipment cancelled", "shipment_id", req.ShipmentID)
	return &contracts.CancelShipmentResponse{Cancelled: true}, nil
}

func main() {
	cfg := config.Load("shipping-service", ":50054")
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
	contracts.RegisterShippingServiceServer(grpcServer, &server{redis: rdb, log: log})
	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()
	log.Info("shipping service listening", "addr", cfg.GRPCAddr)
	if err := grpcServer.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		log.Error("grpc server failed", "error", err)
	}
}
