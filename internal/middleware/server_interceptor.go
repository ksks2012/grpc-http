package middleware

import (
	"context"
	"log"
	"time"

	"github.com/grpc-http/global"
	"github.com/grpc-http/pkg/errcode"
	"github.com/grpc-http/pkg/metatext"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ServerTracing(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}

	parentSpanContext, _ := global.Tracer.Extract(opentracing.TextMap, metatext.MetaDataTextMap{md})
	spanOpts := []opentracing.StartSpanOption{
		opentracing.Tag{Key: string(ext.Component), Value: "gRPC"},
		ext.SpanKindRPCServer,
		ext.RPCServerOption(parentSpanContext),
	}

	span := global.Tracer.StartSpan(info.FullMethod, spanOpts...)
	defer span.Finish()

	ctx = opentracing.ContextWithSpan(ctx, span)
	return handler(ctx, req)
}

func AccessLog(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	requestLog := "access request log: method: %s, begin_time: %d, request: %v"
	beginTime := time.Now().Local().Unix()
	log.Printf(requestLog, info.FullMethod, beginTime, req)

	resp, err := handler(ctx, req)

	responseLog := "access response log: method: %s, begin_time: %d, end_time: %d, response: %v"
	endTime := time.Now().Local().Unix()
	log.Printf(responseLog, info.FullMethod, beginTime, endTime, resp)
	return resp, err
}

func ErrorLog(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)
	if err != nil {
		errLog := "ERROR LOG: method: %s, code: %v, message: %v, details: %v"
		s := errcode.FromError(err)
		log.Printf(errLog, info.FullMethod, s.Code(), s.Err().Error(), s.Details())
	}
	return resp, err
}
