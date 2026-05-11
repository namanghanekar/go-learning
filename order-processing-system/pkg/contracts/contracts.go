package contracts

import (
	"context"

	"google.golang.org/grpc"
)

type OrderItem struct {
	SKU      string `json:"sku"`
	Quantity int32  `json:"quantity"`
}

type CreateOrderRequest struct {
	IdempotencyKey string      `json:"idempotency_key"`
	UserID         string      `json:"user_id"`
	Items          []OrderItem `json:"items"`
	AmountCents    int64       `json:"amount_cents"`
	Email          string      `json:"email"`
}

type CreateOrderResponse struct {
	OrderID string `json:"order_id"`
	UserID  string `json:"user_id,omitempty"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type GetCartRequest struct {
	UserID string `json:"user_id"`
}

type GetCartResponse struct {
	OrderID     string      `json:"order_id"`
	UserID      string      `json:"user_id"`
	Status      string      `json:"status"`
	Items       []OrderItem `json:"items"`
	AmountCents int64       `json:"amount_cents"`
	Email       string      `json:"email"`
	Message     string      `json:"message"`
}

type CheckoutRequest struct {
	IdempotencyKey string `json:"idempotency_key"`
	UserID         string `json:"user_id"`
	PaymentToken   string `json:"payment_token"`
	Address        string `json:"address"`
}

type CheckoutResponse struct {
	OrderID    string `json:"order_id"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	PaymentID  string `json:"payment_id,omitempty"`
	ShipmentID string `json:"shipment_id,omitempty"`
}

type ReserveInventoryRequest struct {
	OrderID string      `json:"order_id"`
	Items   []OrderItem `json:"items"`
}

type ReserveInventoryResponse struct {
	ReservationID string `json:"reservation_id"`
}

type ReleaseInventoryRequest struct {
	ReservationID string `json:"reservation_id"`
}

type ReleaseInventoryResponse struct {
	Released bool `json:"released"`
}

type ProcessPaymentRequest struct {
	OrderID      string `json:"order_id"`
	AmountCents  int64  `json:"amount_cents"`
	PaymentToken string `json:"payment_token"`
}

type ProcessPaymentResponse struct {
	PaymentID string `json:"payment_id"`
}

type RefundPaymentRequest struct {
	PaymentID string `json:"payment_id"`
}

type RefundPaymentResponse struct {
	Refunded bool `json:"refunded"`
}

type CreateShipmentRequest struct {
	OrderID string `json:"order_id"`
	Address string `json:"address"`
}

type CreateShipmentResponse struct {
	ShipmentID string `json:"shipment_id"`
}

type CancelShipmentRequest struct {
	ShipmentID string `json:"shipment_id"`
}

type CancelShipmentResponse struct {
	Cancelled bool `json:"cancelled"`
}

type OrderServiceServer interface {
	CreateOrder(context.Context, *CreateOrderRequest) (*CreateOrderResponse, error)
	GetCart(context.Context, *GetCartRequest) (*GetCartResponse, error)
	Checkout(context.Context, *CheckoutRequest) (*CheckoutResponse, error)
}

type InventoryServiceServer interface {
	ReserveInventory(context.Context, *ReserveInventoryRequest) (*ReserveInventoryResponse, error)
	ReleaseInventory(context.Context, *ReleaseInventoryRequest) (*ReleaseInventoryResponse, error)
}

type PaymentServiceServer interface {
	ProcessPayment(context.Context, *ProcessPaymentRequest) (*ProcessPaymentResponse, error)
	RefundPayment(context.Context, *RefundPaymentRequest) (*RefundPaymentResponse, error)
}

type ShippingServiceServer interface {
	CreateShipment(context.Context, *CreateShipmentRequest) (*CreateShipmentResponse, error)
	CancelShipment(context.Context, *CancelShipmentRequest) (*CancelShipmentResponse, error)
}

func RegisterOrderServiceServer(s grpc.ServiceRegistrar, srv OrderServiceServer) {
	s.RegisterService(&grpc.ServiceDesc{ServiceName: "orders.OrderService", HandlerType: (*OrderServiceServer)(nil), Methods: []grpc.MethodDesc{
		{MethodName: "CreateOrder", Handler: func(i any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
			in := new(CreateOrderRequest)
			if err := dec(in); err != nil {
				return nil, err
			}
			return i.(OrderServiceServer).CreateOrder(ctx, in)
		}},
		{MethodName: "GetCart", Handler: func(i any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
			in := new(GetCartRequest)
			if err := dec(in); err != nil {
				return nil, err
			}
			return i.(OrderServiceServer).GetCart(ctx, in)
		}},
		{MethodName: "Checkout", Handler: func(i any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
			in := new(CheckoutRequest)
			if err := dec(in); err != nil {
				return nil, err
			}
			return i.(OrderServiceServer).Checkout(ctx, in)
		}},
	}}, srv)
}

func RegisterInventoryServiceServer(s grpc.ServiceRegistrar, srv InventoryServiceServer) {
	s.RegisterService(&grpc.ServiceDesc{ServiceName: "orders.InventoryService", HandlerType: (*InventoryServiceServer)(nil), Methods: []grpc.MethodDesc{
		{MethodName: "ReserveInventory", Handler: func(i any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
			in := new(ReserveInventoryRequest)
			if err := dec(in); err != nil {
				return nil, err
			}
			return i.(InventoryServiceServer).ReserveInventory(ctx, in)
		}},
		{MethodName: "ReleaseInventory", Handler: func(i any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
			in := new(ReleaseInventoryRequest)
			if err := dec(in); err != nil {
				return nil, err
			}
			return i.(InventoryServiceServer).ReleaseInventory(ctx, in)
		}},
	}}, srv)
}

func RegisterPaymentServiceServer(s grpc.ServiceRegistrar, srv PaymentServiceServer) {
	s.RegisterService(&grpc.ServiceDesc{ServiceName: "orders.PaymentService", HandlerType: (*PaymentServiceServer)(nil), Methods: []grpc.MethodDesc{
		{MethodName: "ProcessPayment", Handler: func(i any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
			in := new(ProcessPaymentRequest)
			if err := dec(in); err != nil {
				return nil, err
			}
			return i.(PaymentServiceServer).ProcessPayment(ctx, in)
		}},
		{MethodName: "RefundPayment", Handler: func(i any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
			in := new(RefundPaymentRequest)
			if err := dec(in); err != nil {
				return nil, err
			}
			return i.(PaymentServiceServer).RefundPayment(ctx, in)
		}},
	}}, srv)
}

func RegisterShippingServiceServer(s grpc.ServiceRegistrar, srv ShippingServiceServer) {
	s.RegisterService(&grpc.ServiceDesc{ServiceName: "orders.ShippingService", HandlerType: (*ShippingServiceServer)(nil), Methods: []grpc.MethodDesc{
		{MethodName: "CreateShipment", Handler: func(i any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
			in := new(CreateShipmentRequest)
			if err := dec(in); err != nil {
				return nil, err
			}
			return i.(ShippingServiceServer).CreateShipment(ctx, in)
		}},
		{MethodName: "CancelShipment", Handler: func(i any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
			in := new(CancelShipmentRequest)
			if err := dec(in); err != nil {
				return nil, err
			}
			return i.(ShippingServiceServer).CancelShipment(ctx, in)
		}},
	}}, srv)
}

type OrderServiceClient struct{ cc grpc.ClientConnInterface }

func NewOrderServiceClient(cc grpc.ClientConnInterface) *OrderServiceClient {
	return &OrderServiceClient{cc: cc}
}
func (c *OrderServiceClient) CreateOrder(ctx context.Context, in *CreateOrderRequest) (*CreateOrderResponse, error) {
	out := new(CreateOrderResponse)
	err := c.cc.Invoke(ctx, "/orders.OrderService/CreateOrder", in, out)
	return out, err
}
func (c *OrderServiceClient) GetCart(ctx context.Context, in *GetCartRequest) (*GetCartResponse, error) {
	out := new(GetCartResponse)
	err := c.cc.Invoke(ctx, "/orders.OrderService/GetCart", in, out)
	return out, err
}
func (c *OrderServiceClient) Checkout(ctx context.Context, in *CheckoutRequest) (*CheckoutResponse, error) {
	out := new(CheckoutResponse)
	err := c.cc.Invoke(ctx, "/orders.OrderService/Checkout", in, out)
	return out, err
}

type InventoryServiceClient struct{ cc grpc.ClientConnInterface }

func NewInventoryServiceClient(cc grpc.ClientConnInterface) *InventoryServiceClient {
	return &InventoryServiceClient{cc: cc}
}
func (c *InventoryServiceClient) ReserveInventory(ctx context.Context, in *ReserveInventoryRequest) (*ReserveInventoryResponse, error) {
	out := new(ReserveInventoryResponse)
	err := c.cc.Invoke(ctx, "/orders.InventoryService/ReserveInventory", in, out)
	return out, err
}
func (c *InventoryServiceClient) ReleaseInventory(ctx context.Context, in *ReleaseInventoryRequest) (*ReleaseInventoryResponse, error) {
	out := new(ReleaseInventoryResponse)
	err := c.cc.Invoke(ctx, "/orders.InventoryService/ReleaseInventory", in, out)
	return out, err
}

type PaymentServiceClient struct{ cc grpc.ClientConnInterface }

func NewPaymentServiceClient(cc grpc.ClientConnInterface) *PaymentServiceClient {
	return &PaymentServiceClient{cc: cc}
}
func (c *PaymentServiceClient) ProcessPayment(ctx context.Context, in *ProcessPaymentRequest) (*ProcessPaymentResponse, error) {
	out := new(ProcessPaymentResponse)
	err := c.cc.Invoke(ctx, "/orders.PaymentService/ProcessPayment", in, out)
	return out, err
}
func (c *PaymentServiceClient) RefundPayment(ctx context.Context, in *RefundPaymentRequest) (*RefundPaymentResponse, error) {
	out := new(RefundPaymentResponse)
	err := c.cc.Invoke(ctx, "/orders.PaymentService/RefundPayment", in, out)
	return out, err
}

type ShippingServiceClient struct{ cc grpc.ClientConnInterface }

func NewShippingServiceClient(cc grpc.ClientConnInterface) *ShippingServiceClient {
	return &ShippingServiceClient{cc: cc}
}
func (c *ShippingServiceClient) CreateShipment(ctx context.Context, in *CreateShipmentRequest) (*CreateShipmentResponse, error) {
	out := new(CreateShipmentResponse)
	err := c.cc.Invoke(ctx, "/orders.ShippingService/CreateShipment", in, out)
	return out, err
}
func (c *ShippingServiceClient) CancelShipment(ctx context.Context, in *CancelShipmentRequest) (*CancelShipmentResponse, error) {
	out := new(CancelShipmentResponse)
	err := c.cc.Invoke(ctx, "/orders.ShippingService/CancelShipment", in, out)
	return out, err
}
