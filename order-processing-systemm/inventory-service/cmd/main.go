package main

import (
	"fmt"
	"log"
	"net"

	grpcHandler "order-processing-system/inventory-service/internal/grpc"
	"order-processing-system/inventory-service/internal/model"
	"order-processing-system/inventory-service/internal/repository"
	"order-processing-system/inventory-service/internal/service"

	"order-processing-system/shared/postgres"

	pb "order-processing-system/proto/generated"

	grpcServer "google.golang.org/grpc"
)

func main() {

	db := postgres.ConnectDB()

	db.AutoMigrate(&model.Inventory{})

	inventoryRepo := repository.NewInventoryRepository(db)

	inventoryService := service.NewInventoryService(
		inventoryRepo,
	)

	server := grpcHandler.NewInventoryServer(
		inventoryService,
	)

	lis, err := net.Listen(
		"tcp",
		":50051",
	)

	if err != nil {
		log.Fatal(err)
	}

	grpcSrv := grpcServer.NewServer()

	pb.RegisterInventoryServiceServer(
		grpcSrv,
		server,
	)

	fmt.Println(
		"Inventory Service Running On Port 50051",
	)

	if err := grpcSrv.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
