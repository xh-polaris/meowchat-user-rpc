package logic

import (
	"context"

	"github.com/xh-polaris/meowchat-user-rpc/internal/svc"
	"github.com/xh-polaris/meowchat-user-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchUserLogic {
	return &SearchUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchUserLogic) SearchUser(in *pb.SearchUserReq) (*pb.SearchUserResp, error) {
	data, total, err := l.svcCtx.UserModel.SearchUser(l.ctx, in.Nickname, in.Count, in.Skip)
	if err != nil {
		return nil, err
	}
	res := make([]*pb.User, 0, in.Count)
	for _, d := range data {
		m := &pb.User{
			Id:        d.ID.Hex(),
			Nickname:  d.Nickname,
			AvatarUrl: d.AvatarUrl,
		}
		res = append(res, m)
	}
	return &pb.SearchUserResp{Users: res, Total: total}, nil
}
