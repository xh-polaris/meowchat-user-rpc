// Code generated by goctl. DO NOT EDIT.
// Source: user.proto

package user

import (
	"context"

	"github.com/xh-polaris/meowchat-user-rpc/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	GetManyUserReq  = pb.GetManyUserReq
	GetManyUserResp = pb.GetManyUserResp
	GetUserReq      = pb.GetUserReq
	GetUserResp     = pb.GetUserResp
	UpdateUserReq   = pb.UpdateUserReq
	UpdateUserResp  = pb.UpdateUserResp
	UserInfo        = pb.UserInfo

	User interface {
		GetUser(ctx context.Context, in *GetUserReq, opts ...grpc.CallOption) (*GetUserResp, error)
		UpdateUser(ctx context.Context, in *UpdateUserReq, opts ...grpc.CallOption) (*UpdateUserResp, error)
		GetManyUser(ctx context.Context, in *GetManyUserReq, opts ...grpc.CallOption) (*GetManyUserResp, error)
	}

	defaultUser struct {
		cli zrpc.Client
	}
)

func NewUser(cli zrpc.Client) User {
	return &defaultUser{
		cli: cli,
	}
}

func (m *defaultUser) GetUser(ctx context.Context, in *GetUserReq, opts ...grpc.CallOption) (*GetUserResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.GetUser(ctx, in, opts...)
}

func (m *defaultUser) UpdateUser(ctx context.Context, in *UpdateUserReq, opts ...grpc.CallOption) (*UpdateUserResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.UpdateUser(ctx, in, opts...)
}

func (m *defaultUser) GetManyUser(ctx context.Context, in *GetManyUserReq, opts ...grpc.CallOption) (*GetManyUserResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.GetManyUser(ctx, in, opts...)
}
