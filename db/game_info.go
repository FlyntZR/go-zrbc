package db

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// GameInfoDao interface defines CRUD operations for GameInfo
type GameInfoDao interface {
	QueryByID(tx *gorm.DB, gi001 int, gi002 int64, gi003 int) (*GameInfo, error)
	QueryByGameType(tx *gorm.DB, gi001 int) ([]*GameInfo, error)
	QueryByStatus(tx *gorm.DB, gi013 int) ([]*GameInfo, error)
	Create(tx *gorm.DB, gameInfo *GameInfo) error
	Delete(tx *gorm.DB, gi001 int, gi002 int64, gi003 int) error
	Update(tx *gorm.DB, gameInfo *GameInfo) error
	Updates(tx *gorm.DB, gi001 int, gi002 int64, gi003 int, data map[string]interface{}) error
}
type gameInfoDao struct{}

// NewGameInfoDao creates a new instance of GameInfoDao
func NewGameInfoDao() GameInfoDao {
	return &gameInfoDao{}
}

func (dao *gameInfoDao) QueryByID(tx *gorm.DB, gi001 int, gi002 int64, gi003 int) (*GameInfo, error) {
	ret := GameInfo{}
	err := tx.Where("gi001 = ? AND gi002 = ? AND gi003 = ?", gi001, gi002, gi003).First(&ret).Error
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (dao *gameInfoDao) QueryByGameType(tx *gorm.DB, gi001 int) ([]*GameInfo, error) {
	var gameInfos []*GameInfo
	err := tx.Where("gi001 = ?", gi001).Find(&gameInfos).Error
	if err != nil {
		return nil, err
	}
	return gameInfos, nil
}

func (dao *gameInfoDao) QueryByStatus(tx *gorm.DB, gi013 int) ([]*GameInfo, error) {
	var gameInfos []*GameInfo
	err := tx.Where("gi013 = ?", gi013).Find(&gameInfos).Error
	if err != nil {
		return nil, err
	}
	return gameInfos, nil
}

func (dao *gameInfoDao) Create(tx *gorm.DB, gameInfo *GameInfo) error {
	return tx.Table("game_info").Create(gameInfo).Error
}

func (dao *gameInfoDao) Delete(tx *gorm.DB, gi001 int, gi002 int64, gi003 int) error {
	return tx.Where("gi001 = ? AND gi002 = ? AND gi003 = ?", gi001, gi002, gi003).Delete(&GameInfo{}).Error
}

func (dao *gameInfoDao) Update(tx *gorm.DB, gameInfo *GameInfo) error {
	return tx.Table("game_info").Where("gi001 = ? AND gi002 = ? AND gi003 = ?",
		gameInfo.Gi001, gameInfo.Gi002, gameInfo.Gi003).Updates(gameInfo).Error
}

func (dao *gameInfoDao) Updates(tx *gorm.DB, gi001 int, gi002 int64, gi003 int, data map[string]interface{}) error {
	return tx.Table("game_info").Where("gi001 = ? AND gi002 = ? AND gi003 = ?",
		gi001, gi002, gi003).Updates(data).Error
}

// GameInfo represents the game_info table structure
type GameInfo struct {
	Gi001   int             `db:"gi001"`   // 遊戲類別
	Gi002   int64           `db:"gi002"`   // 場次編號
	Gi003   int             `db:"gi003"`   // 子場次編號
	Gi004   time.Time       `db:"gi004"`   // 開局時間
	Gi005   decimal.Decimal `db:"gi005"`   // 總下注金額
	Gi006   time.Time       `db:"gi006"`   // 開獎時間
	Gi007   string          `db:"gi007"`   // 開獎內容
	Gi008   decimal.Decimal `db:"gi008"`   // 輸贏
	Gi009   decimal.Decimal `db:"gi009"`   // 總退水
	Gi010   decimal.Decimal `db:"gi010"`   // 結果：贏得＋退水－下注
	Gi011   int             `db:"gi011"`   // Table ID
	Gi012   int             `db:"gi012"`   // result
	Gi013   int             `db:"gi013"`   // 狀態：0:未開,1:己開,2:重對,9:取消
	Eid     int             `db:"eid"`     // 荷官
	Serid   int             `db:"serid"`   // 服務編號
	ThirdNo string          `db:"3rdno"`   // 第三方編號
	IsLock  string          `db:"is_lock"` // 鎖定
}

// TableName returns the table name
func (g *GameInfo) TableName() string {
	return "game_info"
}
