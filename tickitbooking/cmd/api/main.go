package main

import (
	"context"
	"log"
	"net/http"
	"worldtour-tickets/internal/api"
	"worldtour-tickets/internal/config"
	"worldtour-tickets/internal/contracts"
	"worldtour-tickets/internal/rpcjson"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding"
)

func main() {
	cfg := config.Load()
	encoding.RegisterCodec(rpcjson.Codec{})

	inventoryConn, err := grpc.NewClient(
		cfg.InventoryGRPCAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.ForceCodec(rpcjson.Codec{})),
	)
	if err != nil {
		log.Fatalf("connect inventory: %v", err)
	}
	defer inventoryConn.Close()

	paymentConn, err := grpc.NewClient(
		cfg.PaymentGRPCAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.ForceCodec(rpcjson.Codec{})),
	)
	if err != nil {
		log.Fatalf("connect payment: %v", err)
	}
	defer paymentConn.Close()

	server, err := api.NewServer(
		contracts.NewInventoryServiceClient(inventoryConn),
		contracts.NewPaymentServiceClient(paymentConn),
		cfg.SeatRoomID,
	)
	if err != nil {
		log.Fatalf("build api: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server.StartInventoryListener(ctx)

	log.Printf("api listening on %s", cfg.APIServerAddr)
	if err := http.ListenAndServe(cfg.APIServerAddr, server.Handler()); err != nil {
		log.Fatalf("serve api: %v", err)
	}
}
