package config

import "os"

type Config struct {
	InventoryGRPCAddr string
	PaymentGRPCAddr   string
	APIServerAddr     string
	MySQLDSN          string
	SeatRoomID        string
}

func Load() Config {
	return Config{
		InventoryGRPCAddr: envOrDefault("INVENTORY_GRPC_ADDR", ":50051"),
		PaymentGRPCAddr:   envOrDefault("PAYMENT_GRPC_ADDR", ":50052"),
		APIServerAddr:     envOrDefault("API_SERVER_ADDR", ":8080"),
		MySQLDSN:          envOrDefault("MYSQL_DSN", "root:naman@tcp(127.0.0.1:3306)/ticketbooking?charset=utf8mb4&parseTime=True&loc=Local"),
		SeatRoomID:        envOrDefault("DEFAULT_ROOM_ID", "world-tour-main-stage"),
	}
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
