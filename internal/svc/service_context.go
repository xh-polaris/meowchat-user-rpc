package svc

import (
	"github.com/xh-polaris/meowchat-user-rpc/internal/config"
	"github.com/xh-polaris/meowchat-user-rpc/internal/model"
)

type ServiceContext struct {
	Config config.Config
	model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUserModel(c.Mongo.URL, c.Mongo.DB, model.UserCollectionName, c.CacheConf, c.Elasticsearch),
	}
}
