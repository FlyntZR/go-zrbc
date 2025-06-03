package service

import (
	"context"
	"go-zrbc/db"
	"go-zrbc/pkg/xlog"
	"go-zrbc/service"
	"go-zrbc/view"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-redis/redis/v8"

	"gorm.io/gorm"
)

const (
	TokenPrefix = "zrys:access_token"
)

type UserService interface {
	//用户信息
	GetUserInfo(ctx context.Context, userID int64) (*view.GetUserInfoResp, error)
	GetUserByAccountAndPwd(ctx context.Context, account, pwd string) (*view.GetUserInfoResp, error)
}

type userService struct {
	userDao db.UserDao

	s3Client *s3.Client
	redisCli *redis.Client
	*service.Session
}

func NewUserService(
	sess *service.Session,
	userDao db.UserDao,

	s3Client *s3.Client,
	redisCli *redis.Client,
) UserService {
	srv := &userService{
		redisCli: redisCli,
		userDao:  userDao,
		s3Client: s3Client,
	}
	srv.Session = sess
	return srv
}

func (srv *userService) GetUserInfo(ctx context.Context, userID int64) (*view.GetUserInfoResp, error) {
	var ret *db.Member
	var err error
	err = srv.Tx(func(tx *gorm.DB) error {
		ret, err = srv.userDao.QueryByID(tx, userID)
		if err != nil {
			xlog.Errorf("error to get user info, err:%+v", err)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	resp := view.GetUserInfoResp{
		User: DBToViewUser(ret),
	}
	return &resp, nil
}

func (srv *userService) GetUserByAccountAndPwd(ctx context.Context, account, pwd string) (*view.GetUserInfoResp, error) {
	var ret *db.Member
	var err error
	err = srv.Tx(func(tx *gorm.DB) error {
		ret, err = srv.userDao.QueryByAccountAndPwd(tx, account, pwd)
		if err != nil {
			xlog.Errorf("error to get user, err:%+v", err)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	resp := view.GetUserInfoResp{
		User: DBToViewUser(ret),
	}
	return &resp, nil
}
