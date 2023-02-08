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

// GetManyUser
// 若查询失败不返回异常，得到空结构体/*
func (l *GetManyUserLogic) GetManyUser(in *pb.GetManyUserReq) (*pb.GetManyUserResp, error) {
	resp := make([]*pb.UserInfo, 0, len(in.UserId))
	for _, userid := range in.UserId {
		res, err := l.svcCtx.UserModel.FindOne(l.ctx, userid)
		if err != nil {
			resp = append(resp, &pb.UserInfo{})
		} else {
			resp = append(resp, &pb.UserInfo{
				UserId:    res.ID.Hex(),
				AvatarUrl: res.AvatarUrl,
				Nickname:  res.Nickname,
			})
		}
	}
	return &pb.GetManyUserResp{UserInfo: resp}, nil
}
