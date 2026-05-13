# =========================
# DOCKER
# =========================

up:
	docker-compose up -d

down:
	docker-compose down

restart:
	docker-compose restart

logs:
	docker-compose logs -f


# =========================
# ORDER SERVICE
# =========================

order:
	cd order-service && go run cmd/main.go


# =========================
# INVENTORY SERVICE
# =========================

inventory:
	cd inventory-service && go run cmd/main.go


# =========================
# PAYMENT SERVICE
# =========================

payment:
	cd payment-service && go run cmd/main.go


# =========================
# SHIPMENT SERVICE
# =========================

shipment:
	cd shipment-service && go run cmd/main.go


# =========================
# NOTIFICATION SERVICE
# =========================

notification:
	cd notification-service && go run cmd/main.go


# =========================
# API GATEWAY
# =========================

gateway:
	cd api-gateway && go run cmd/main.go


# =========================
# RUN ALL SERVICES
# =========================

run-all:
	start cmd /k "cd order-service && go run cmd/main.go"
	start cmd /k "cd inventory-service && go run cmd/main.go"
	start cmd /k "cd payment-service && go run cmd/main.go"
	start cmd /k "cd shipment-service && go run cmd/main.go"
	start cmd /k "cd notification-service && go run cmd/main.go"
	start cmd /k "cd api-gateway && go run cmd/main.go"


# =========================
# GO MODULES
# =========================

tidy:
	cd order-service && go mod tidy
	cd inventory-service && go mod tidy
	cd payment-service && go mod tidy
	cd shipment-service && go mod tidy
	cd notification-service && go mod tidy
	cd api-gateway && go mod tidy


# =========================
# PROTO GENERATION
# =========================

proto:
	protoc --go_out=. --go-grpc_out=. proto/inventory.proto
	protoc --go_out=. --go-grpc_out=. proto/payment.proto


# =========================
# KAFKA TOPIC
# =========================

topic:
	docker exec kafka kafka-topics --create --topic order.created --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1