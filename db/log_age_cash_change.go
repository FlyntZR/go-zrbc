package db

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type LogAgeCashChangeDao interface {
	QueryByID(tx *gorm.DB, id int64) (*LogAgeCashChange, error)
	QueryByLvAndID(tx *gorm.DB, lv int64, id int64) (*LogAgeCashChange, error)
	Create(tx *gorm.DB, record *LogAgeCashChange) (int64, error)
	DeleteByID(tx *gorm.DB, id int64) error
	Update(tx *gorm.DB, record *LogAgeCashChange) error
	Updates(tx *gorm.DB, id int64, data map[string]interface{}) error
	QueryByOpID(tx *gorm.DB, opID int64) ([]*LogAgeCashChange, error)
}

type logAgeCashChangeDao struct{}

func NewLogAgeCashChangeDao() LogAgeCashChangeDao {
	return &logAgeCashChangeDao{}
}

func (dao *logAgeCashChangeDao) QueryByID(tx *gorm.DB, id int64) (*LogAgeCashChange, error) {
	ret := LogAgeCashChange{}
	err := tx.First(&ret, id).Error
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (dao *logAgeCashChangeDao) QueryByLvAndID(tx *gorm.DB, lv int64, id int64) (*LogAgeCashChange, error) {
	ret := LogAgeCashChange{}
	err := tx.Where("lacc02 = ? AND lacc03 = ?", lv, id).First(&ret).Error
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (dao *logAgeCashChangeDao) Create(tx *gorm.DB, record *LogAgeCashChange) (int64, error) {
	err := tx.Table(TableNameLogAgeCashChange).Create(record).Error
	if err != nil {
		return 0, err
	}
	return record.Lacc01, nil
}

func (dao *logAgeCashChangeDao) DeleteByID(tx *gorm.DB, id int64) error {
	return tx.Where("lacc01 = ?", id).Delete(&LogAgeCashChange{}).Error
}

func (dao *logAgeCashChangeDao) Update(tx *gorm.DB, record *LogAgeCashChange) error {
	return tx.Table(TableNameLogAgeCashChange).Where("lacc01 = ?", record.Lacc01).Updates(record).Error
}

func (dao *logAgeCashChangeDao) Updates(tx *gorm.DB, id int64, data map[string]interface{}) error {
	return tx.Table(TableNameLogAgeCashChange).Where("lacc01 = ?", id).Updates(data).Error
}

func (dao *logAgeCashChangeDao) QueryByOpID(tx *gorm.DB, opID int64) ([]*LogAgeCashChange, error) {
	var records []*LogAgeCashChange
	err := tx.Where("lacc05 = ?", opID).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

const TableNameLogAgeCashChange = "log_age_cash_change"

// LogAgeCashChange mapped from table <log_age_cash_change>
type LogAgeCashChange struct {
	Lacc01    int64           `gorm:"column:lacc01;primaryKey;autoIncrement:true" json:"lacc01"`
	Lacc02    int             `gorm:"column:lacc02;not null;comment:LV" json:"lacc02"`               // LV
	Lacc03    int64           `gorm:"column:lacc03;not null;comment:ID" json:"lacc03"`               // ID
	OpUtp     string          `gorm:"column:op_utp;not null;comment:操作者類型" json:"op_utp"`            // 操作者類型
	Lacc04    int64           `gorm:"column:lacc04;not null;comment:OPLV" json:"lacc04"`             // OPLV
	Lacc05    int64           `gorm:"column:lacc05;not null;comment:OPID" json:"lacc05"`             // OPID
	Lacc06    decimal.Decimal `gorm:"column:lacc06;not null;comment:異動金額" json:"lacc06"`             // 異動金額
	Lacc07    time.Time       `gorm:"column:lacc07;not null;comment:操作時間" json:"lacc07"`             // 操作時間
	Lacc08    string          `gorm:"column:lacc08;not null;comment:備註" json:"lacc08"`               // 備註
	Lacc09    decimal.Decimal `gorm:"column:lacc09;not null;comment:異動前金額" json:"lacc09"`            // 異動前金額
	Lacc10    int64           `gorm:"column:lacc10;not null;comment:upid" json:"lacc10"`             // upid
	Lacc11    decimal.Decimal `gorm:"column:lacc11;not null;comment:上層異動前金額" json:"lacc11"`          // 上層異動前金額
	Pointtype int             `gorm:"column:pointtype;not null;comment:0:現金;1:信用;" json:"pointtype"` // 0:現金;1:信用;
}

// TableName LogAgeCashChange's table name
func (*LogAgeCashChange) TableName() string {
	return TableNameLogAgeCashChange
}
