package model

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/mitchellh/mapstructure"
	"github.com/xh-polaris/meowchat-user-rpc/internal/config"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/monc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
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
		SearchUser(ctx context.Context, name string, skip, count int64) ([]*User, int64, error)
	}

	customUserModel struct {
		*defaultUserModel
		es        *elasticsearch.Client
		indexName string
	}
)

// NewUserModel returns a model for the mongo.
func NewUserModel(url, db, collection string, c cache.CacheConf, es config.ElasticsearchConf) UserModel {
	conn := monc.MustNewModel(url, db, collection, c)
	esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: es.Addresses,
		Username:  es.Username,
		Password:  es.Password,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	return &customUserModel{
		defaultUserModel: newDefaultUserModel(conn),
		es:               esClient,
		indexName:        fmt.Sprintf("%s.%s-alias", db, UserCollectionName),
	}
}

func (m customUserModel) UpsertUser(ctx context.Context, data *User) error {
	key := prefixUserCacheKey + data.ID.Hex()

	filter := bson.M{
		ID: data.ID,
	}

	set := bson.M{
		UpdateAt: time.Now(),
	}
	if data.Nickname != "" {
		set[Nickname] = data.Nickname
	}
	if data.AvatarUrl != "" {
		set[AvatarUrl] = data.AvatarUrl
	}

	update := bson.M{
		"$set": set,
		"$setOnInsert": bson.M{
			ID:       data.ID,
			CreateAt: time.Now(),
		},
	}

	option := options.UpdateOptions{}
	option.SetUpsert(true)

	_, err := m.conn.UpdateOne(ctx, key, filter, update, &option)
	return err
}

func (m customUserModel) SearchUser(ctx context.Context, name string, count, skip int64) ([]*User, int64, error) {
	search := m.es.Search
	query := map[string]any{
		"from": skip,
		"size": count,
		"query": map[string]any{
			"bool": map[string]any{
				"must": []any{
					map[string]any{
						"multi_match": map[string]any{
							"query":  name,
							"fields": []string{Nickname},
						},
					},
				},
			},
		},
		"sort": map[string]any{
			"_score": map[string]any{
				"order": "desc",
			},
			CreateAt: map[string]any{
				"order": "desc",
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, 0, err
	}
	res, err := search(
		search.WithIndex(m.indexName),
		search.WithContext(ctx),
		search.WithBody(&buf),
	)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, 0, err
		} else {
			logx.Errorf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}
	var r map[string]any
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, 0, err
	}
	hits := r["hits"].(map[string]any)["hits"].([]any)
	total := int64(r["hits"].(map[string]any)["total"].(map[string]any)["value"].(float64))
	users := make([]*User, 0, 10)
	for i := range hits {
		hit := hits[i].(map[string]any)
		user := &User{}
		source := hit["_source"].(map[string]any)
		if source[CreateAt], err = time.Parse("2006-01-02T15:04:05Z07:00", source[CreateAt].(string)); err != nil {
			return nil, 0, err
		}
		if source[UpdateAt], err = time.Parse("2006-01-02T15:04:05Z07:00", source[UpdateAt].(string)); err != nil {
			return nil, 0, err
		}
		hit["_source"] = source
		err := mapstructure.Decode(hit["_source"], user)
		if err != nil {
			return nil, 0, err
		}
		oid := hit[ID].(string)
		id, err := primitive.ObjectIDFromHex(oid)
		if err != nil {
			return nil, 0, err
		}
		user.ID = id
		users = append(users, user)
	}
	return users, total, nil
}
