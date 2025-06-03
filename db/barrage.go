package db

import (
	"time"

	"gorm.io/gorm"
)

type BarrageDao interface {
	CreateBarrage(tx *gorm.DB, barrage *Barrage) (int64, error)
	DeleteByID(tx *gorm.DB, uniqueID int64) error
	GetBarragesByVideoSeriesID(tx *gorm.DB, ps, pn int, videoSeriesID int64) ([]*Barrage, error)
	GetBarragesByVideoSeriesIDAndPlaySeconds(tx *gorm.DB, ps, pn int, videoSeriesID int64, playSeconds int) ([]*Barrage, error)
}

type barrageDao struct{}

func NewBarrageDao() BarrageDao {
	return &barrageDao{}
}

func (dao *barrageDao) CreateBarrage(tx *gorm.DB, barrage *Barrage) (int64, error) {
	return barrage.ID, tx.Table("ys_barrage").Create(barrage).Error
}

func (dao *barrageDao) DeleteByID(tx *gorm.DB, uniqueID int64) error {
	return tx.Where("id", uniqueID).Delete(&Barrage{}).Error
}

func (dao *barrageDao) GetBarragesByVideoSeriesID(tx *gorm.DB, ps, pn int, videoSeriesID int64) ([]*Barrage, error) {
	ret := []*Barrage{}
	if err := tx.Table("barrage").Where("video_series_id=?", videoSeriesID).Order("play_seconds").Order("created_at desc").Offset(ps * (pn - 1)).Limit(ps).Find(&ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

func (dao *barrageDao) GetBarragesByVideoSeriesIDAndPlaySeconds(tx *gorm.DB, ps, pn int, videoSeriesID int64, playSeconds int) ([]*Barrage, error) {
	ret := []*Barrage{}
	if err := tx.Table("barrage").Where("video_series_id=? and play_seconds>=?", videoSeriesID, playSeconds).Order("created_at").Offset(ps * (pn - 1)).Limit(ps).Find(&ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

const TableNameBarrage = "barrage"

// Barrage 弹幕表
type Barrage struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement:true;comment:id(自增)" json:"id"`                    // id(自增)
	VideoSeriesID int64     `gorm:"column:video_series_id;not null;comment:剧集id" json:"video_series_id"`                 // 剧集id
	MemberID      int64     `gorm:"column:member_id;comment:会员id" json:"member_id"`                                      // 会员id
	DeviceID      string    `gorm:"column:device_id;comment:浏览器指纹（设备id）" json:"device_id"`                               // 浏览器指纹（设备id）
	Content       string    `gorm:"column:content;not null;comment:弹幕内容" json:"content"`                                 // 弹幕内容
	PlaySeconds   int       `gorm:"column:play_seconds;comment:已播放秒数" json:"play_seconds"`                               // 已播放秒数
	CreatedAt     time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
}

// TableName Barrage's table name
func (*Barrage) TableName() string {
	return TableNameBarrage
}
