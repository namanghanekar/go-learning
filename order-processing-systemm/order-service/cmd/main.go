package main

import (
	"fmt"
	"log"
	"net"

	"order-processing-system/order-service/internal/client"
	grpcHandler "order-processing-system/order-service/internal/grpc"
	"order-processing-system/order-service/internal/kafka"
	"order-processing-system/order-service/internal/model"
	"order-processing-system/order-service/internal/repository"
	"order-processing-system/order-service/internal/service"

	"order-processing-system/shared/postgres"

	pb "order-processing-system/proto/generated"

	grpcServer "google.golang.org/grpc"
)

func main() {

	db := postgres.ConnectDB()

	db.AutoMigrate(&model.Order{})

	orderRepo := repository.NewOrderRepository(db)

	inventoryClient := client.NewInventoryClient()

	paymentClient := client.NewPaymentClient()

	producer := kafka.NewProducer()

	orderService := service.NewOrderService(
		orderRepo,
		inventoryClient,
		paymentClient,
		producer,
	)

	server := grpcHandler.NewOrderServer(
		orderService,
	)

	lis, err := net.Listen(
		"tcp",
		":50050",
	)

	if err != nil {
		log.Fatal(err)
	}

	grpcSrv := grpcServer.NewServer()

	pb.RegisterOrderServiceServer(
		grpcSrv,
		server,
	)

	fmt.Println(
		"Order Service Running On Port 50050",
	)

	if err := grpcSrv.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
