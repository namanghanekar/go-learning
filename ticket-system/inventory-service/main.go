package main

import (
	"log"
	"net"

	"ticket-system/inventory-service/models"
	"ticket-system/inventory-service/repository"
	"ticket-system/inventory-service/service"
	"ticket-system/shared"

	"google.golang.org/grpc"

	grpcserver "ticket-system/inventory-service/grpc"
	pb "ticket-system/proto"
)

func main() {
	shared.InitDB()
	shared.DB.AutoMigrate(&models.Seat{})
	repo := repository.NewRepo(shared.DB)
	svc := service.NewService(repo)
	if err := svc.SeedSeats(100); err != nil {
		log.Fatal(err)
	}
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	pb.RegisterInventoryServiceServer(s, grpcserver.NewServer(svc))
	log.Println("inventory service running on :50051")
	s.Serve(lis)
}
