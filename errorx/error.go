package errorx

import (
	"github.com/xh-polaris/meowchat-user-rpc/internal/model"
	"google.golang.org/grpc/status"
)

var (
	ErrNotFound        = status.Error(12001, "data not found")
	ErrInvalidObjectId = status.Error(12002, "invalid objectId")
)

func Switch(err error) error {
	switch err {
	case model.ErrNotFound:
		return ErrNotFound
	case model.ErrInvalidObjectId:
		return ErrInvalidObjectId
	default:
		return err
	}
}
