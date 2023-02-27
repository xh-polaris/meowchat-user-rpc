package logic

import (
	"context"
	"github.com/xh-polaris/meowchat-user-rpc/errorx"
	"github.com/xh-polaris/meowchat-user-rpc/internal/model"
	"github.com/xh-polaris/meowchat-user-rpc/internal/svc"
	"github.com/xh-polaris/meowchat-user-rpc/pb"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserLogic) UpdateUser(in *pb.UpdateUserReq) (*pb.UpdateUserResp, error) {
	oid, err := primitive.ObjectIDFromHex(in.User.Id)
	if err != nil {
		return nil, errorx.ErrInvalidObjectId
	}

	err = l.svcCtx.UserModel.UpsertUser(l.ctx, &model.User{
		ID:        oid,
		AvatarUrl: in.User.AvatarUrl,
		Nickname:  in.User.Nickname,
	})
	if err != nil {
		return nil, err
	}

	return &pb.UpdateUserResp{}, nil
}
