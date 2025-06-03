package service

import (
	"context"
	"errors"
	"fmt"
	"go-zrbc/db"
	"go-zrbc/pkg/utils"
	"go-zrbc/pkg/xlog"
	"go-zrbc/service"
	"go-zrbc/view"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type WebService interface {
	GetTimeTs(ctx context.Context) (*view.GetTimeTsResp, error)

	CreateBarrage(ctx context.Context, req *view.Barrage, userID int64) error
	GetBarragesByVideoSeriesID(ctx context.Context, ps, pn int, videoSeriesID int64) (*view.GetBarragesResp, error)
	GetBarragesByVideoSeriesIDAndPlaySeconds(ctx context.Context, ps, pn int, videoSeriesID int64, playSeconds int) (*view.GetBarragesResp, error)
	DeleteBarrage(ctx context.Context, req *view.DeleteBarrageReq) error
}

type webService struct {
	barrageDao db.BarrageDao
	s3Client   *s3.Client
	redisCli   *redis.Client

	*service.Session
}

func NewWebService(
	sess *service.Session,
	barrageDao db.BarrageDao,
	s3Client *s3.Client,
	redisCli *redis.Client,
) WebService {
	srv := &webService{
		Session:    sess,
		barrageDao: barrageDao,
		s3Client:   s3Client,
		redisCli:   redisCli,
	}
	return srv
}

func (srv *webService) GetTimeTs(ctx context.Context) (*view.GetTimeTsResp, error) {
	loc, _ := time.LoadLocation("Local")
	now := time.Now()
	return &view.GetTimeTsResp{
		STm: fmt.Sprintf("current time: %s", now),
		Tm:  fmt.Sprintf("current time: %s", now),
		Ts:  time.Now().Unix(),
		Tz:  loc.String(),
	}, nil
}

func (srv *webService) CreateBarrage(ctx context.Context, req *view.Barrage, userID int64) error {
	if userID <= 0 && req.DeviceID == "" {
		return errors.New("client info err")
	}
	videoSeriesID, err := utils.ParseInt(req.VideoSeriesID)
	if err != nil {
		return err
	}
	barrage := db.Barrage{
		VideoSeriesID: videoSeriesID,
		MemberID:      userID,
		DeviceID:      req.DeviceID,
		Content:       req.Content,
		PlaySeconds:   req.PlaySeconds,
		CreatedAt:     time.Now(),
	}

	return srv.Tx(func(tx *gorm.DB) error {
		_, err := srv.barrageDao.CreateBarrage(tx, &barrage)
		if err != nil {
			xlog.Errorf("error to create barrage, err:%+v", err)
			return err
		}
		return nil
	})
}

func (srv *webService) GetBarragesByVideoSeriesID(ctx context.Context, ps, pn int, videoSeriesID int64) (*view.GetBarragesResp, error) {
	ret := []*db.Barrage{}
	var err error
	err = srv.Tx(func(tx *gorm.DB) error {
		ret, err = srv.barrageDao.GetBarragesByVideoSeriesID(tx, ps, pn, videoSeriesID)
		if err != nil {
			xlog.Errorf("error to get barrages, err:%+v", err)
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	resp := view.GetBarragesResp{
		YsBarrages: make([]*view.BarrageResp, 0),
		Total:      0,
	}
	for i := 0; i < len(ret); i++ {
		tv := DBToViewBarrage(ret[i])
		resp.YsBarrages = append(resp.YsBarrages, tv)
	}
	resp.Total = int64(len(resp.YsBarrages))
	return &resp, nil
}

func (srv *webService) GetBarragesByVideoSeriesIDAndPlaySeconds(ctx context.Context, ps, pn int, videoSeriesID int64, playSeconds int) (*view.GetBarragesResp, error) {
	ret := []*db.Barrage{}
	var err error
	err = srv.Tx(func(tx *gorm.DB) error {
		ret, err = srv.barrageDao.GetBarragesByVideoSeriesIDAndPlaySeconds(tx, ps, pn, videoSeriesID, playSeconds)
		if err != nil {
			xlog.Errorf("error to get barrages, err:%+v", err)
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	resp := view.GetBarragesResp{
		YsBarrages: make([]*view.BarrageResp, 0),
		Total:      0,
	}
	for i := 0; i < len(ret); i++ {
		tv := DBToViewBarrage(ret[i])
		resp.YsBarrages = append(resp.YsBarrages, tv)
	}
	resp.Total = int64(len(resp.YsBarrages))
	return &resp, nil
}

func (srv *webService) DeleteBarrage(ctx context.Context, req *view.DeleteBarrageReq) error {
	return srv.Tx(func(tx *gorm.DB) error {
		err := srv.barrageDao.DeleteByID(tx, int64(req.ID))
		if err != nil {
			xlog.Errorf("error to delete barrage by id, id:%d, err:%+v", int64(req.ID), err)
			return err
		}
		return nil
	})
}
