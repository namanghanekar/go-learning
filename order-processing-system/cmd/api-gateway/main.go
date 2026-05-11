package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"

	"order-processing-system/internal/config"
	"order-processing-system/internal/grpcx"
	"order-processing-system/internal/logx"
	"order-processing-system/internal/shutdown"
	"order-processing-system/pkg/contracts"
)

type gateway struct {
	orders *contracts.OrderServiceClient
	log    *slog.Logger
}

func (g gateway) cart(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		g.addToCart(w, r)
	case http.MethodGet:
		g.getCart(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (g gateway) addToCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req contracts.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	resp, err := g.orders.CreateOrder(ctx, &req)
	if err != nil {
		g.log.Error("cart save failed", "error", err)
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (g gateway) getCart(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id query parameter is required", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	resp, err := g.orders.GetCart(ctx, &contracts.GetCartRequest{UserID: userID})
	if err != nil {
		g.log.Error("cart fetch failed", "error", err)
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (g gateway) checkout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req contracts.CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	resp, err := g.orders.Checkout(ctx, &req)
	if err != nil {
		g.log.Error("checkout failed", "error", err)
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func main() {
	cfg := config.Load("api-gateway", ":0")
	log := logx.New(cfg.ServiceName)
	ctx := shutdown.Context()
	orderAddr := env("ORDER_GRPC_ADDR", "localhost:50051")
	conn, err := grpcx.Dial(ctx, orderAddr)
	if err != nil {
		log.Error("order service dial failed", "error", err)
		return
	}
	defer conn.Close()
	g := gateway{orders: contracts.NewOrderServiceClient(conn), log: log}
	mux := http.NewServeMux()
	mux.HandleFunc("/cart", g.cart)
	mux.HandleFunc("/checkout", g.checkout)
	mux.HandleFunc("/orders", g.addToCart)
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("ok")) })
	server := &http.Server{Addr: cfg.HTTPAddr, Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()
	log.Info("api gateway listening", "addr", cfg.HTTPAddr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("http server failed", "error", err)
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
