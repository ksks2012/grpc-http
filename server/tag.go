package server

import (
	"context"
	"encoding/json"

	// "net/http"

	bapi "github.com/grpc-http/pkg/bapi"
	"github.com/grpc-http/pkg/errcode"
	pb "github.com/grpc-http/proto"
	"google.golang.org/grpc/metadata"
)

type TagServer struct {
	pb.TagServiceServer
	auth *Auth
}

type Auth struct{}

func (a *Auth) GetAppKey() string {
	// TODO: get from DB
	return "hong"
}

func (a *Auth) GetAppSecret() string {
	// TODO: get from DB
	return "blog-service"
}

func (a *Auth) Check(ctx context.Context) error {
	md, _ := metadata.FromIncomingContext(ctx)

	var appKey, appSecret string
	if value, ok := md["app_key"]; ok {
		appKey = value[0]
	}
	if value, ok := md["app_secret"]; ok {
		appSecret = value[0]
	}
	if appKey != a.GetAppKey() || appSecret != a.GetAppSecret() {
		return errcode.TogRPCError(errcode.Unauthorized)
	}

	return nil
}

func NewTagServer() *TagServer {
	return &TagServer{}
}

func (t *TagServer) GetTagList(ctx context.Context, r *pb.GetTagListRequest) (*pb.GetTagListReply, error) {
	if err := t.auth.Check(ctx); err != nil {
		return nil, err
	}
	api := bapi.NewAPI("http://127.0.0.1:18080")
	body, err := api.GetTagList(ctx, r.GetName())
	if err != nil {
		return nil, errcode.TogRPCError(errcode.ErrorGetTagListFail)
	}

	tagList := pb.GetTagListReply{}
	err = json.Unmarshal(body, &tagList)
	if err != nil {
		return nil, errcode.TogRPCError(errcode.Fail)
	}

	return &tagList, nil
}
