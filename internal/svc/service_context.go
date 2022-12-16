package svc

import (
	"github.com/xh-polaris/meowchat-user-rpc/internal/config"
	"github.com/xh-polaris/meowchat-user-rpc/internal/model"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config config.Config
	model.UserModel
	*redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUserModel(c.Mongo.URL, c.Mongo.DB, model.UserCollectionName, c.CacheConf),
		Redis:     c.Redis.NewRedis(),
	}
}
