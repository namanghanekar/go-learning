package main

import (
	"context"
	"log"
	"net"
	"worldtour-tickets/internal/config"
	"worldtour-tickets/internal/contracts"
	"worldtour-tickets/internal/payment"
	"worldtour-tickets/internal/rpcjson"
	"worldtour-tickets/internal/store"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding"
)

func main() {
	cfg := config.Load()
	encoding.RegisterCodec(rpcjson.Codec{})

	db, err := store.OpenMySQL(cfg.MySQLDSN)
	if err != nil {
		log.Fatalf("open payment mysql: %v", err)
	}

	repo := store.NewPaymentRepository(db)

	inventoryConn, err := grpc.NewClient(
		cfg.InventoryGRPCAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.ForceCodec(rpcjson.Codec{})),
	)
	if err != nil {
		log.Fatalf("connect inventory: %v", err)
	}
	defer inventoryConn.Close()

	paymentService := payment.NewService(repo, contracts.NewInventoryServiceClient(inventoryConn))

	appCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	paymentService.Start(appCtx)

	grpcServer := grpc.NewServer(grpc.ForceServerCodec(rpcjson.Codec{}))
	contracts.RegisterPaymentServiceServer(grpcServer, paymentService)

	listener, err := net.Listen("tcp", cfg.PaymentGRPCAddr)
	if err != nil {
		log.Fatalf("listen payment: %v", err)
	}

	log.Printf("payment listening on %s", cfg.PaymentGRPCAddr)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("serve payment: %v", err)
	}
}
