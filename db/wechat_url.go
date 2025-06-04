package db

import (
	"errors"

	"gorm.io/gorm"
)

type WechatURLDao interface {
	Create(tx *gorm.DB, wechatURL *WechatURL) (int64, error)
	QueryByID(tx *gorm.DB, id int64) (*WechatURL, error)
	QueryByURL(tx *gorm.DB, url string) (*WechatURL, error)
	UpdateWechatURL(tx *gorm.DB, wechatURL *WechatURL) error
	UpdatesWechatURL(tx *gorm.DB, id int64, data map[string]interface{}) error
	DeleteByID(tx *gorm.DB, id int64) error
	ListWechatURLs(tx *gorm.DB) ([]*WechatURL, error)
	GetRandomWechatURL(tx *gorm.DB) (*WechatURL, error)
	UpdateWechatURLUseCount(tx *gorm.DB, id int64) error
}

type wechatURLDao struct{}

func NewWechatURLDao() WechatURLDao {
	return &wechatURLDao{}
}

func (dao *wechatURLDao) Create(tx *gorm.DB, wechatURL *WechatURL) (int64, error) {
	err := tx.Table(TableNameWechatURL).Create(wechatURL).Error
	if err != nil {
		return 0, err
	}
	return wechatURL.ID, nil
}

func (dao *wechatURLDao) QueryByID(tx *gorm.DB, id int64) (*WechatURL, error) {
	ret := WechatURL{}
	err := tx.First(&ret, id).Error
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (dao *wechatURLDao) QueryByURL(tx *gorm.DB, url string) (*WechatURL, error) {
	ret := WechatURL{}
	err := tx.Where("url = ?", url).First(&ret).Error
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (dao *wechatURLDao) UpdateWechatURL(tx *gorm.DB, wechatURL *WechatURL) error {
	return tx.Table(TableNameWechatURL).Where("id = ?", wechatURL.ID).Updates(wechatURL).Error
}

func (dao *wechatURLDao) UpdatesWechatURL(tx *gorm.DB, id int64, data map[string]interface{}) error {
	return tx.Table(TableNameWechatURL).Where("id = ?", id).Updates(data).Error
}

func (dao *wechatURLDao) DeleteByID(tx *gorm.DB, id int64) error {
	return tx.Where("id = ?", id).Delete(&WechatURL{}).Error
}

func (dao *wechatURLDao) ListWechatURLs(tx *gorm.DB) ([]*WechatURL, error) {
	var urls []*WechatURL
	err := tx.Find(&urls).Error
	if err != nil {
		return nil, err
	}
	return urls, nil
}

// GetRandomWechatURL retrieves a random Wechat URL
func (dao *wechatURLDao) GetRandomWechatURL(tx *gorm.DB) (*WechatURL, error) {
	var urlData WechatURL
	err := tx.Order(gorm.Expr("random()")).First(&urlData).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no Wechat URLs available")
		}
		return nil, err
	}
	return &urlData, nil
}

// UpdateWechatURLUseCount increments the use count for a Wechat URL
func (dao *wechatURLDao) UpdateWechatURLUseCount(tx *gorm.DB, id int64) error {
	return tx.Table(TableNameWechatURL).
		Where("id = ?", id).
		UpdateColumn("usecount", gorm.Expr("usecount + 1")).
		Error
}

const TableNameWechatURL = "wechat_url"

// WechatURL mapped from table <wechat_url>
type WechatURL struct {
	ID       int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	URL      string `gorm:"column:url;not null" json:"url"`
	Status   int    `gorm:"column:status;not null;comment:0:使用,1:禁用" json:"status"` // 0:使用,1:禁用
	UseCount int    `gorm:"column:usecount;not null;comment:使用次數" json:"usecount"`  // 使用次數
}

// TableName WechatURL's table name
func (*WechatURL) TableName() string {
	return TableNameWechatURL
}
