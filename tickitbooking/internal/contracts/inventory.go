package contracts

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

type ListSeatsRequest struct {
	RoomID string `json:"roomId"`
}

type ListSeatsResponse struct {
	Seats []SeatDTO `json:"seats"`
}

type HoldSeatRequest struct {
	RoomID     string `json:"roomId"`
	SeatNumber string `json:"seatNumber"`
	UserID     string `json:"userId"`
	HoldToken  string `json:"holdToken"`
	TTLSeconds int32  `json:"ttlSeconds"`
}

type HoldSeatResponse struct {
	Seat SeatDTO `json:"seat"`
}

type ConfirmSeatRequest struct {
	SeatID     string `json:"seatId"`
	HoldToken  string `json:"holdToken"`
	PaymentRef string `json:"paymentRef"`
}

type ReleaseSeatRequest struct {
	SeatID    string `json:"seatId"`
	HoldToken string `json:"holdToken"`
	Reason    string `json:"reason"`
}

type SeatCommandAck struct {
	Accepted bool      `json:"accepted"`
	Message  string    `json:"message"`
	Seat     SeatDTO   `json:"seat"`
	At       time.Time `json:"at"`
}

type StreamEnvelope struct {
	ClientName string              `json:"clientName"`
	Event      *SeatEvent          `json:"event,omitempty"`
	Confirm    *ConfirmSeatRequest `json:"confirm,omitempty"`
	Release    *ReleaseSeatRequest `json:"release,omitempty"`
	Ack        *SeatCommandAck     `json:"ack,omitempty"`
}

type InventoryServiceClient interface {
	ListSeats(ctx context.Context, in *ListSeatsRequest, opts ...grpc.CallOption) (*ListSeatsResponse, error)
	HoldSeat(ctx context.Context, in *HoldSeatRequest, opts ...grpc.CallOption) (*HoldSeatResponse, error)
	OpenSeatStream(ctx context.Context, opts ...grpc.CallOption) (InventoryService_OpenSeatStreamClient, error)
}

type InventoryServiceServer interface {
	ListSeats(context.Context, *ListSeatsRequest) (*ListSeatsResponse, error)
	HoldSeat(context.Context, *HoldSeatRequest) (*HoldSeatResponse, error)
	OpenSeatStream(InventoryService_OpenSeatStreamServer) error
}

type InventoryService_OpenSeatStreamClient interface {
	Send(*StreamEnvelope) error
	Recv() (*StreamEnvelope, error)
	grpc.ClientStream
}

type InventoryService_OpenSeatStreamServer interface {
	Send(*StreamEnvelope) error
	Recv() (*StreamEnvelope, error)
	grpc.ServerStream
}

type inventoryServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewInventoryServiceClient(cc grpc.ClientConnInterface) InventoryServiceClient {
	return &inventoryServiceClient{cc: cc}
}

func (c *inventoryServiceClient) ListSeats(ctx context.Context, in *ListSeatsRequest, opts ...grpc.CallOption) (*ListSeatsResponse, error) {
	out := new(ListSeatsResponse)
	err := c.cc.Invoke(ctx, "/inventory.InventoryService/ListSeats", in, out, opts...)
	return out, err
}

func (c *inventoryServiceClient) HoldSeat(ctx context.Context, in *HoldSeatRequest, opts ...grpc.CallOption) (*HoldSeatResponse, error) {
	out := new(HoldSeatResponse)
	err := c.cc.Invoke(ctx, "/inventory.InventoryService/HoldSeat", in, out, opts...)
	return out, err
}

func (c *inventoryServiceClient) OpenSeatStream(ctx context.Context, opts ...grpc.CallOption) (InventoryService_OpenSeatStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &InventoryService_ServiceDesc.Streams[0], "/inventory.InventoryService/OpenSeatStream", opts...)
	if err != nil {
		return nil, err
	}
	return &inventoryServiceOpenSeatStreamClient{ClientStream: stream}, nil
}

type inventoryServiceOpenSeatStreamClient struct {
	grpc.ClientStream
}

func (x *inventoryServiceOpenSeatStreamClient) Send(m *StreamEnvelope) error {
	return x.ClientStream.SendMsg(m)
}
func (x *inventoryServiceOpenSeatStreamClient) Recv() (*StreamEnvelope, error) {
	m := new(StreamEnvelope)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func RegisterInventoryServiceServer(s grpc.ServiceRegistrar, srv InventoryServiceServer) {
	s.RegisterService(&InventoryService_ServiceDesc, srv)
}

var InventoryService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "inventory.InventoryService",
	HandlerType: (*InventoryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListSeats",
			Handler:    _InventoryService_ListSeats_Handler,
		},
		{
			MethodName: "HoldSeat",
			Handler:    _InventoryService_HoldSeat_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "OpenSeatStream",
			Handler:       _InventoryService_OpenSeatStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
}

func _InventoryService_ListSeats_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(ListSeatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InventoryServiceServer).ListSeats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/inventory.InventoryService/ListSeats"}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(InventoryServiceServer).ListSeats(ctx, req.(*ListSeatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _InventoryService_HoldSeat_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(HoldSeatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InventoryServiceServer).HoldSeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/inventory.InventoryService/HoldSeat"}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(InventoryServiceServer).HoldSeat(ctx, req.(*HoldSeatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _InventoryService_OpenSeatStream_Handler(srv any, stream grpc.ServerStream) error {
	return srv.(InventoryServiceServer).OpenSeatStream(&inventoryServiceOpenSeatStreamServer{ServerStream: stream})
}

type inventoryServiceOpenSeatStreamServer struct {
	grpc.ServerStream
}

func (x *inventoryServiceOpenSeatStreamServer) Send(m *StreamEnvelope) error {
	return x.ServerStream.SendMsg(m)
}
func (x *inventoryServiceOpenSeatStreamServer) Recv() (*StreamEnvelope, error) {
	m := new(StreamEnvelope)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}
