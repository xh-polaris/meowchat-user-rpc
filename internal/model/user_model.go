package model

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/monc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const UserCollectionName = "user"

var _ UserModel = (*customUserModel)(nil)

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
		UpsertUser(ctx context.Context, data *User) error
	}

	customUserModel struct {
		*defaultUserModel
	}
)

func (m customUserModel) UpsertUser(ctx context.Context, data *User) error {
	key := prefixUserCacheKey + data.ID.Hex()

	filter := bson.M{
		"_id": data.ID,
	}

	set := bson.M{
		"updateAt": time.Now(),
	}
	if data.Nickname != "" {
		set["nickname"] = data.Nickname
	}
	if data.AvatarUrl != "" {
		set["avatarUrl"] = data.AvatarUrl
	}

	update := bson.M{
		"$set": set,
		"$setOnInsert": bson.M{
			"_id":      data.ID,
			"createAt": time.Now(),
		},
	}

	option := options.UpdateOptions{}
	option.SetUpsert(true)

	_, err := m.conn.UpdateOne(ctx, key, filter, update, &option)
	return err
}

// NewUserModel returns a model for the mongo.
func NewUserModel(url, db, collection string, c cache.CacheConf) UserModel {
	conn := monc.MustNewModel(url, db, collection, c)
	return &customUserModel{
		defaultUserModel: newDefaultUserModel(conn),
	}
}
