package main

import (
	"fmt"
	"log"
	"net"

	grpcHandler "order-processing-system/payment-service/internal/grpc"

	"order-processing-system/payment-service/internal/model"
	"order-processing-system/payment-service/internal/repository"
	"order-processing-system/payment-service/internal/service"

	"order-processing-system/shared/postgres"

	pb "order-processing-system/proto/generated"

	grpcServer "google.golang.org/grpc"
)

func main() {

	db := postgres.ConnectDB()

	db.AutoMigrate(&model.Payment{})

	paymentRepo := repository.NewPaymentRepository(db)

	paymentService := service.NewPaymentService(
		paymentRepo,
	)

	server := grpcHandler.NewPaymentServer(
		paymentService,
	)

	lis, err := net.Listen(
		"tcp",
		":50052",
	)

	if err != nil {
		log.Fatal(err)
	}

	grpcSrv := grpcServer.NewServer()

	pb.RegisterPaymentServiceServer(
		grpcSrv,
		server,
	)

	fmt.Println(
		"Payment Service Running On Port 50052",
	)

	if err := grpcSrv.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
