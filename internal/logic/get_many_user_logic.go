package logic

import (
	"context"

	"github.com/xh-polaris/meowchat-user-rpc/internal/svc"
	"github.com/xh-polaris/meowchat-user-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetManyUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetManyUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetManyUserLogic {
	return &GetManyUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetManyUserLogic) GetManyUser(in *pb.GetManyUserReq) (*pb.GetManyUserResp, error) {
	resp := make([]*pb.UserInfo, len(in.UserId))
	for _, userid := range in.UserId {
		res, err := l.svcCtx.UserModel.FindOne(l.ctx, userid)
		if err != nil {
			return nil, err
		}
		resp = append(resp, &pb.UserInfo{
			UserId:    res.ID.Hex(),
			AvatarUrl: res.AvatarUrl,
			Nickname:  res.Nickname,
		})
	}
	return &pb.GetManyUserResp{UserInfo: resp}, nil
}
