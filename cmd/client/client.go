package main

import (
	"context"
	"log"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/grpc-http/global"
	"github.com/grpc-http/internal/middleware"
	"github.com/grpc-http/pkg/tracer"
	pb "github.com/grpc-http/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func init() {
	err := setupTracer()
	if err != nil {
		log.Fatalf("init.setupTracer err: %v", err)
	}
}

type Auth struct {
	AppKey    string
	AppSecret string
}

func (a *Auth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"app_key": a.AppKey, "app_secret": a.AppSecret}, nil
}

func (a *Auth) RequireTransportSecurity() bool {
	return false
}

func main() {
	auth := Auth{
		AppKey:    "hong",
		AppSecret: "blog-service",
	}
	ctx := context.Background()
	opts := []grpc.DialOption{grpc.WithPerRPCCredentials(&auth)}
	// FIXME: target
	clientConn, err := GetClientConn(ctx, "localhost:8004", opts)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	defer clientConn.Close()

	tagServiceClient := pb.NewTagServiceClient(clientConn)
	resp, _ := tagServiceClient.GetTagList(ctx, &pb.GetTagListRequest{Name: "Go"})

	log.Printf("resp: %v", resp)
}

func GetClientConn(ctx context.Context, target string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			middleware.UnaryContextTimeout(),
		),
	))
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
	opts = append(opts, grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			middleware.ClientTracing(),
		),
	))

	return grpc.DialContext(ctx, target, opts...)
}

func setupTracer() error {
	jaegerTracer, _, err := tracer.NewJaegerTracer("grpc-http client", "127.0.0.1:6831")
	if err != nil {
		return err
	}
	global.Tracer = jaegerTracer
	return nil
}
