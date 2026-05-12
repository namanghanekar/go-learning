package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"strconv"
	"time"

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
	contracts.InventoryServiceServer
	redis *redis.Client
	log   *slog.Logger
}

type inventoryDecrement struct {
	key string
	qty int64
}

func (s *server) ReserveInventory(ctx context.Context, req *contracts.ReserveInventoryRequest) (*contracts.ReserveInventoryResponse, error) {
	reservationID := "res_" + uuid.NewString()
	var rolledBack []inventoryDecrement
	for _, item := range req.Items {
		if item.Quantity <= 0 {
			s.rollbackDecrements(ctx, rolledBack)
			return nil, errors.New("quantity must be positive")
		}
		key := "inventory:" + item.SKU
		qty := int64(item.Quantity)
		remaining, err := s.redis.DecrBy(ctx, key, qty).Result()
		if err != nil {
			s.rollbackDecrements(ctx, rolledBack)
			return nil, err
		}
		if remaining < 0 {
			_, _ = s.redis.IncrBy(ctx, key, qty).Result()
			s.rollbackDecrements(ctx, rolledBack)
			return nil, errors.New("insufficient inventory for sku " + item.SKU)
		}
		rolledBack = append(rolledBack, inventoryDecrement{key: key, qty: qty})
	}
	tx := s.redis.TxPipeline()
	for _, item := range req.Items {
		tx.HSet(ctx, "reservation:"+reservationID, item.SKU, item.Quantity)
	}
	tx.Expire(ctx, "reservation:"+reservationID, 24*time.Hour)
	if _, err := tx.Exec(ctx); err != nil {
		s.rollbackDecrements(ctx, rolledBack)
		return nil, err
	}
	s.log.Info("inventory reserved", "order_id", req.OrderID, "reservation_id", reservationID)
	return &contracts.ReserveInventoryResponse{ReservationID: reservationID}, nil
}

func (s *server) rollbackDecrements(ctx context.Context, ops []inventoryDecrement) {
	for i := len(ops) - 1; i >= 0; i-- {
		_, _ = s.redis.IncrBy(ctx, ops[i].key, ops[i].qty).Result()
	}
}

func (s *server) ReleaseInventory(ctx context.Context, req *contracts.ReleaseInventoryRequest) (*contracts.ReleaseInventoryResponse, error) {
	values, err := s.redis.HGetAll(ctx, "reservation:"+req.ReservationID).Result()
	if err != nil {
		return nil, err
	}
	for sku, qty := range values {
		n, err := strconv.ParseInt(qty, 10, 64)
		if err != nil {
			continue
		}
		_, _ = s.redis.IncrBy(ctx, "inventory:"+sku, n).Result()
	}
	_ = s.redis.Del(ctx, "reservation:"+req.ReservationID).Err()
	s.log.Info("inventory released", "reservation_id", req.ReservationID)
	return &contracts.ReleaseInventoryResponse{Released: true}, nil
}

func main() {
	cfg := config.Load("inventory-service", ":50052")
	log := logx.New(cfg.ServiceName)
	ctx := shutdown.Context()
	rdb, err := redisx.Connect(ctx, cfg.RedisAddr)
	if err != nil {
		log.Error("redis connection failed", "error", err)
		return
	}
	for _, sku := range []string{"SKU-BOOK", "SKU-PHONE", "SKU-LAPTOP"} {
		_ = rdb.SetNX(ctx, "inventory:"+sku, 100, 0).Err()
	}
	listener, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		log.Error("listen failed", "error", err)
		return
	}
	grpcServer := grpcx.Server()
	contracts.RegisterInventoryServiceServer(grpcServer, &server{redis: rdb, log: log})
	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()
	log.Info("inventory service listening", "addr", cfg.GRPCAddr)
	if err := grpcServer.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		log.Error("grpc server failed", "error", err)
	}
}
