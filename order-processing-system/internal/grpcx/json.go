package grpcx

import (
	"context"
	"encoding/json"

	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
)

const CodecName = "json"

type JSONCodec struct{}

func (JSONCodec) Name() string                       { return CodecName }
func (JSONCodec) Marshal(v any) ([]byte, error)      { return json.Marshal(v) }
func (JSONCodec) Unmarshal(data []byte, v any) error { return json.Unmarshal(data, v) }

func init() {
	encoding.RegisterCodec(JSONCodec{})
}

func Dial(ctx context.Context, target string) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, target, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.ForceCodec(JSONCodec{})))
}

func Server() *grpc.Server {
	return grpc.NewServer(grpc.ForceServerCodec(JSONCodec{}))
}
