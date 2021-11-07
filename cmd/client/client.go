package main

import (
	"context"
	"log"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	pb "github.com/grpc-http/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func main() {
	ctx := context.Background()
	clientConn, _ := GetClientConn(
		ctx,
		"localhost:18080",
		[]grpc.DialOption{grpc.WithBlock()},
	)
	defer clientConn.Close()

	tagServiceClient := pb.NewTagServiceClient(clientConn)
	resp, _ := tagServiceClient.GetTagList(ctx, &pb.GetTagListRequest{Name: "Go"})

	log.Printf("resp: %v", resp)
}

func GetClientConn(ctx context.Context, target string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	// opts []grpc.DialOption
	opts = append(opts, grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			grpc_retry.UnaryClientInterceptor(
				grpc_retry.WithMax(2),
				grpc_retry.WithCodes(
					codes.Unknown,
					codes.Internal,
					codes.DeadlineExceeded,
				),
			),
		),
	))
	return grpc.DialContext(ctx, target, opts...)
}
