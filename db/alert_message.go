package db

import (
	"time"

	"gorm.io/gorm"
)

// AlertMessageDao interface defines CRUD operations for AlertMessage
type AlertMessageDao interface {
	QueryByID(tx *gorm.DB, id int64) (*AlertMessage, error)
	QueryByMid(tx *gorm.DB, mid int64) ([]*AlertMessage, error)
	Create(tx *gorm.DB, alertMessage *AlertMessage) (int64, error)
	DeleteByID(tx *gorm.DB, id int64) error
	Update(tx *gorm.DB, alertMessage *AlertMessage) error
	Updates(tx *gorm.DB, id int64, data map[string]interface{}) error
	QueryUnhandledMessages(tx *gorm.DB) ([]*AlertMessage, error)
}

type alertMessageDao struct{}

// NewAlertMessageDao creates a new instance of AlertMessageDao
func NewAlertMessageDao() AlertMessageDao {
	return &alertMessageDao{}
}

func (dao *alertMessageDao) QueryByID(tx *gorm.DB, id int64) (*AlertMessage, error) {
	ret := AlertMessage{}
	err := tx.First(&ret, id).Error
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (dao *alertMessageDao) QueryByMid(tx *gorm.DB, mid int64) ([]*AlertMessage, error) {
	var messages []*AlertMessage
	err := tx.Where("mid = ?", mid).Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (dao *alertMessageDao) Create(tx *gorm.DB, alertMessage *AlertMessage) (int64, error) {
	err := tx.Table(TableNameAlertMessage).Create(alertMessage).Error
	if err != nil {
		return 0, err
	}
	return alertMessage.ID, nil
}

func (dao *alertMessageDao) DeleteByID(tx *gorm.DB, id int64) error {
	return tx.Where("id = ?", id).Delete(&AlertMessage{}).Error
}

func (dao *alertMessageDao) Update(tx *gorm.DB, alertMessage *AlertMessage) error {
	return tx.Table(TableNameAlertMessage).Where("id = ?", alertMessage.ID).Updates(alertMessage).Error
}

func (dao *alertMessageDao) Updates(tx *gorm.DB, id int64, data map[string]interface{}) error {
	return tx.Table(TableNameAlertMessage).Where("id = ?", id).Updates(data).Error
}

func (dao *alertMessageDao) QueryUnhandledMessages(tx *gorm.DB) ([]*AlertMessage, error) {
	var messages []*AlertMessage
	err := tx.Where("status = ?", 0).Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}

const TableNameAlertMessage = "alertMessage"

// AlertMessage mapped from table <alertMessage>
type AlertMessage struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Mid          int64     `gorm:"column:mid;not null" json:"mid"`
	Message      string    `gorm:"column:message;not null" json:"message"`
	ErrorTime    time.Time `gorm:"column:errorTime;not null" json:"errorTime"`
	UnierrorTime int64     `gorm:"column:unierrorTime;not null" json:"unierrorTime"`
	Operator     string    `gorm:"column:operator;not null;comment:操作者" json:"operator"`    // 操作者
	Status       int       `gorm:"column:status;not null;comment:0:未處理1:已處理" json:"status"` // 0:未處理1:已處理
}

// TableName AlertMessage's table name
func (*AlertMessage) TableName() string {
	return TableNameAlertMessage
}
