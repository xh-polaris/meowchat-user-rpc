package logic

import (
	"context"

	"github.com/xh-polaris/meowchat-user-rpc/errorx"
	"github.com/xh-polaris/meowchat-user-rpc/internal/svc"
	"github.com/xh-polaris/meowchat-user-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserLogic) GetUser(in *pb.GetUserReq) (*pb.GetUserResp, error) {
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, in.UserId)
	if err != nil {
		return nil, errorx.Switch(err)
	}

	return &pb.GetUserResp{
		User: &pb.User{
			Id:        user.ID.Hex(),
			AvatarUrl: user.AvatarUrl,
			Nickname:  user.Nickname,
		},
	}, nil
}
