package db

import (
	"gorm.io/gorm"
)

type BetLimitDefaultDao interface {
	QueryAll(tx *gorm.DB) ([]*BetLimitDefault, error)
	QueryByID(tx *gorm.DB, id int64) (*BetLimitDefault, error)
	QueryByGtype(tx *gorm.DB, gtype int64) ([]*BetLimitDefault, error)
	Create(tx *gorm.DB, betLimit *BetLimitDefault) (int64, error)
	DeleteByID(tx *gorm.DB, id int64) error
	Update(tx *gorm.DB, betLimit *BetLimitDefault) error
	Updates(tx *gorm.DB, id int64, data map[string]interface{}) error
}

type betLimitDefaultDao struct{}

func NewBetLimitDefaultDao() BetLimitDefaultDao {
	return &betLimitDefaultDao{}
}

func (dao *betLimitDefaultDao) QueryAll(tx *gorm.DB) ([]*BetLimitDefault, error) {
	ret := []*BetLimitDefault{}
	err := tx.Find(&ret).Error
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (dao *betLimitDefaultDao) QueryByID(tx *gorm.DB, id int64) (*BetLimitDefault, error) {
	ret := BetLimitDefault{}
	err := tx.First(&ret, id).Error
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (dao *betLimitDefaultDao) QueryByGtype(tx *gorm.DB, gtype int64) ([]*BetLimitDefault, error) {
	var ret []*BetLimitDefault
	err := tx.Where("gtype = ?", gtype).Find(&ret).Error
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (dao *betLimitDefaultDao) Create(tx *gorm.DB, betLimit *BetLimitDefault) (int64, error) {
	err := tx.Table(TableNameBetLimitDefault).Create(betLimit).Error
	if err != nil {
		return 0, err
	}
	return betLimit.ID, nil
}

func (dao *betLimitDefaultDao) DeleteByID(tx *gorm.DB, id int64) error {
	return tx.Where("id = ?", id).Delete(&BetLimitDefault{}).Error
}

func (dao *betLimitDefaultDao) Update(tx *gorm.DB, betLimit *BetLimitDefault) error {
	return tx.Table(TableNameBetLimitDefault).Where("id = ?", betLimit.ID).Updates(betLimit).Error
}

func (dao *betLimitDefaultDao) Updates(tx *gorm.DB, id int64, data map[string]interface{}) error {
	return tx.Table(TableNameBetLimitDefault).Where("id = ?", id).Updates(data).Error
}

const TableNameBetLimitDefault = "bet_limit_default"

// BetLimitDefault mapped from table <bet_limit_default>
type BetLimitDefault struct {
	ID     int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Gtype  int64  `gorm:"column:gtype;not null" json:"gtype"`
	Set01  string `gorm:"column:set01;not null" json:"set01"`
	Set02  string `gorm:"column:set02;not null" json:"set02"`
	Set03  string `gorm:"column:set03;not null" json:"set03"`
	Set04  string `gorm:"column:set04;not null" json:"set04"`
	Set05  string `gorm:"column:set05;not null" json:"set05"`
	Set06  string `gorm:"column:set06;not null" json:"set06"`
	Set07  string `gorm:"column:set07;not null" json:"set07"`
	Set08  string `gorm:"column:set08;not null" json:"set08"`
	Set09  string `gorm:"column:set09;not null" json:"set09"`
	Set10  string `gorm:"column:set10;not null" json:"set10"`
	Set11  string `gorm:"column:set11;not null" json:"set11"`
	Set12  string `gorm:"column:set12;not null" json:"set12"`
	Set13  string `gorm:"column:set13;not null" json:"set13"`
	Set14  string `gorm:"column:set14;not null" json:"set14"`
	Status int64  `gorm:"column:status;not null;default:1" json:"status"`
	Sort   int64  `gorm:"column:sort;not null;comment:排序" json:"sort"` // 排序
}

// TableName BetLimitDefault's table name
func (*BetLimitDefault) TableName() string {
	return TableNameBetLimitDefault
}
