package db

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

const TableNameGameType = "game_type"

type GameTypeDao interface {
	QueryByID(tx *gorm.DB, code int64) (*GameType, error)
	QueryByStatus(tx *gorm.DB, status int) ([]*GameType, error)
	CreateGameType(tx *gorm.DB, gameType *GameType) (int64, error)
	DeleteByCode(tx *gorm.DB, code int64) error
	UpdateGameType(tx *gorm.DB, gameType *GameType) error
	UpdatesGameType(tx *gorm.DB, code int64, data map[string]interface{}) error
}

type gameTypeDao struct{}

func NewGameTypeDao() GameTypeDao {
	return &gameTypeDao{}
}

func (dao *gameTypeDao) QueryByID(tx *gorm.DB, code int64) (*GameType, error) {
	ret := GameType{}
	err := tx.First(&ret, code).Error
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (dao *gameTypeDao) QueryByStatus(tx *gorm.DB, status int) ([]*GameType, error) {
	var gameTypes []*GameType
	err := tx.Where("status = ?", status).Find(&gameTypes).Error
	if err != nil {
		return nil, err
	}
	return gameTypes, nil
}

func (dao *gameTypeDao) CreateGameType(tx *gorm.DB, gameType *GameType) (int64, error) {
	err := tx.Table(TableNameGameType).Create(gameType).Error
	if err != nil {
		return 0, err
	}
	return gameType.Code, nil
}

func (dao *gameTypeDao) DeleteByCode(tx *gorm.DB, code int64) error {
	return tx.Where("code = ?", code).Delete(&GameType{}).Error
}

func (dao *gameTypeDao) UpdateGameType(tx *gorm.DB, gameType *GameType) error {
	return tx.Table(TableNameGameType).Where("code = ?", gameType.Code).Updates(gameType).Error
}

func (dao *gameTypeDao) UpdatesGameType(tx *gorm.DB, code int64, data map[string]interface{}) error {
	return tx.Table(TableNameGameType).Where("code = ?", code).Updates(data).Error
}

// GameType mapped from table <game_type>
type GameType struct {
	Code   int64           `gorm:"column:code;primaryKey" json:"code"`
	Twname string          `gorm:"column:twname;not null" json:"twname"`
	Cnname string          `gorm:"column:cnname;not null" json:"cnname"`
	Enname string          `gorm:"column:enname;not null" json:"enname"`
	Inname string          `gorm:"column:inname;not null" json:"inname"`
	Janame string          `gorm:"column:janame;not null" json:"janame"`
	Koname string          `gorm:"column:koname;not null" json:"koname"`
	Thname string          `gorm:"column:thname;not null" json:"thname"`
	Viname string          `gorm:"column:viname;not null" json:"viname"`
	Hiname string          `gorm:"column:hiname;not null" json:"hiname"`
	Msname string          `gorm:"column:msname;not null" json:"msname"`
	Esname string          `gorm:"column:esname;not null" json:"esname"`
	Status int             `gorm:"column:status;not null;default:1" json:"status"`
	Ps003  decimal.Decimal `gorm:"column:ps003;not null" json:"ps003"`
	Ps004  int             `gorm:"column:ps004;not null" json:"ps004"`
	Ps005  int             `gorm:"column:ps005;not null" json:"ps005"`
	Ps006  int             `gorm:"column:ps006;not null" json:"ps006"`
	Ps007  int             `gorm:"column:ps007;not null" json:"ps007"`
	Ps008  int             `gorm:"column:ps008;not null" json:"ps008"`
	Ps009  int             `gorm:"column:ps009;not null" json:"ps009"`
	Ps010  int             `gorm:"column:ps010;not null" json:"ps010"`
	Ps011  int             `gorm:"column:ps011;not null" json:"ps011"`
}

// TableName GameType's table name
func (*GameType) TableName() string {
	return TableNameGameType
}
