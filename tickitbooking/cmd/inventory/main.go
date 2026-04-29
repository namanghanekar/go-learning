package main

import (
	"context"
	"log"
	"net"
	"worldtour-tickets/internal/config"
	"worldtour-tickets/internal/contracts"
	"worldtour-tickets/internal/inventory"
	"worldtour-tickets/internal/rpcjson"
	"worldtour-tickets/internal/store"

	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
)

func main() {
	cfg := config.Load()
	encoding.RegisterCodec(rpcjson.Codec{})

	db, err := store.OpenMySQL(cfg.MySQLDSN)
	if err != nil {
		log.Fatalf("open inventory mysql: %v", err)
	}

	seatRepo := store.NewSeatRepository(db)
	if err := seatRepo.SeedRoom(cfg.SeatRoomID, 20); err != nil {
		log.Fatalf("seed seats: %v", err)
	}

	inventoryService := inventory.NewService(seatRepo)
	inventoryService.StartReconciler(context.Background(), cfg.SeatRoomID)

	grpcServer := grpc.NewServer(grpc.ForceServerCodec(rpcjson.Codec{}))
	contracts.RegisterInventoryServiceServer(grpcServer, inventoryService)

	listener, err := net.Listen("tcp", cfg.InventoryGRPCAddr)
	if err != nil {
		log.Fatalf("listen inventory: %v", err)
	}

	log.Printf("inventory listening on %s", cfg.InventoryGRPCAddr)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("serve inventory: %v", err)
	}
}
