package service

import (
	"go-zrbc/db"
	"go-zrbc/view"
)

func DBToViewBarrage(dbBarrage *db.Barrage) *view.BarrageResp {
	return &view.BarrageResp{
		ID:          dbBarrage.ID,
		MemberID:    dbBarrage.MemberID,
		DeviceID:    dbBarrage.DeviceID,
		Content:     dbBarrage.Content,
		PlaySeconds: dbBarrage.PlaySeconds,
		CreatedAt:   dbBarrage.CreatedAt.UnixMilli(),
	}
}
