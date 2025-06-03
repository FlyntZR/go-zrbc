package token

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go-zrbc/pkg/xlog"

	"github.com/go-redis/redis/v8"
)

// swagger:model
type TokenUser struct {
	ID       int64  `json:"id"`
	UserName string `json:"name"`
	Mobile   string `json:"mobile"`
}

var ErrTokenIsExpired = errors.New("token is expired")

func GetToken(redisCli *redis.Client, tPrefix, token string) (*TokenUser, error) {
	tKey := fmt.Sprintf("%s:%s", tPrefix, token)
	data, err := redisCli.Get(context.TODO(), tKey).Result()
	if err != nil {
		if err == redis.Nil {
			xlog.Errorf("error get token, tkey:%s, err:%+v", tKey, err)
			return nil, ErrTokenIsExpired
		}
		xlog.Error(err)
		return nil, err
	}
	user := TokenUser{}
	err = json.Unmarshal([]byte(data), &user)
	return &user, err
}

func DelToken(redisCli *redis.Client, tPrefix, token string) error {
	tKey := fmt.Sprintf("%s:%s", tPrefix, token)
	_, err := redisCli.Del(context.TODO(), tKey).Result()
	if err != nil {
		xlog.Error(err)
	}
	return err
}

func SetToken(redisCli *redis.Client, user *TokenUser, tPrefix, token string) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	tKey := fmt.Sprintf("%s:%s", tPrefix, token)
	_, err = redisCli.Set(context.TODO(), tKey, data, time.Hour*24).Result()
	if err != nil {
		xlog.Error(err)
	}
	return err
}
