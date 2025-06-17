package db

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

const TableNameBet01 = "bet01"

// Bet01Dao interface defines CRUD operations for Bet01
type Bet01Dao interface {
	QueryByID(tx *gorm.DB, id int64) (*Bet01, error)
	Create(tx *gorm.DB, bet01 *Bet01) (int64, error)
	DeleteByID(tx *gorm.DB, id int64) error
	Update(tx *gorm.DB, bet01 *Bet01) error
	Updates(tx *gorm.DB, id int64, data map[string]interface{}) error
	GetBet01ListForUnsettledReport(tx *gorm.DB, date time.Time) ([]*Bet01Summary, error)
}

type bet01Dao struct{}

func NewBet01Dao() Bet01Dao {
	return &bet01Dao{}
}

func (dao *bet01Dao) QueryByID(tx *gorm.DB, id int64) (*Bet01, error) {
	ret := Bet01{}
	err := tx.First(&ret, id).Error
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (dao *bet01Dao) Create(tx *gorm.DB, bet01 *Bet01) (int64, error) {
	err := tx.Table(TableNameBet01).Create(bet01).Error
	if err != nil {
		return 0, err
	}
	return bet01.Bet01, nil
}

func (dao *bet01Dao) DeleteByID(tx *gorm.DB, id int64) error {
	return tx.Where("bet01 = ?", id).Delete(&Bet01{}).Error
}

func (dao *bet01Dao) Update(tx *gorm.DB, bet01 *Bet01) error {
	return tx.Table(TableNameBet01).Where("bet01 = ?", bet01.Bet01).Updates(bet01).Error
}

func (dao *bet01Dao) Updates(tx *gorm.DB, id int64, data map[string]interface{}) error {
	return tx.Table(TableNameBet01).Where("bet01 = ?", id).Updates(data).Error
}

type Bet01Summary struct {
	BetID      int64           `gorm:"column:betId" json:"betId"`
	GID        int             `gorm:"column:gid" json:"gid"`
	Event      decimal.Decimal `gorm:"column:event" json:"event"`
	EventChild int             `gorm:"column:eventChild" json:"eventChild"`
	ID         int             `gorm:"column:id" json:"id"`
	BetTime    time.Time       `gorm:"column:betTime" json:"betTime"`
	BetResult  string          `gorm:"column:betResult" json:"betResult"`
	Bet        decimal.Decimal `gorm:"column:bet" json:"bet"`
	AID        int             `gorm:"column:aid" json:"aid"`
	TableID    int             `gorm:"column:tableId" json:"tableId"`
	Commission int             `gorm:"column:commission" json:"commission"`
	Round      decimal.Decimal `gorm:"column:round" json:"round"`
	SubRound   int             `gorm:"column:subround" json:"subround"`
	GName      string          `gorm:"column:gname" json:"gname"`
	User       string          `gorm:"column:user" json:"user"`
}

func (dao *bet01Dao) GetBet01ListForUnsettledReport(tx *gorm.DB, date time.Time) ([]*Bet01Summary, error) {
	var ret []*Bet01Summary

	err := tx.Table(TableNameBet01).
		Joins("LEFT JOIN game_type ON game_type.Code = bet02").Joins("LEFT JOIN member ON bet05 = member.mem001").
		Select("bet01 as betId, bet02 as gid, bet03 as event, bet04 as eventChild, bet05 as id, "+
			"bet08 as betTime, bet09 as betResult, bet13 as bet, bet19 as aid, bet31 as tableId, "+
			"commission, bet03 as round, bet04 as subround, game_type.cnname as gname, member.mem002 as user").
		Where("bet01 NOT IN (SELECT bet01 FROM bet02)").
		Where("bet07 = ?", date.Format("2006-01-02")).
		Where("bet02 NOT LIKE ?", "301").
		Where("bet30 = ?", "N").
		Find(&ret).Error

	if err != nil {
		return nil, err
	}
	return ret, nil
}

// Bet01 mapped from table <bet01>
type Bet01 struct {
	Bet01        int64           `gorm:"column:bet01;primaryKey;comment:注單編號" json:"bet01"`                    // 注單編號
	Bet02        int             `gorm:"column:bet02;not null;comment:遊戲類別編號" json:"bet02"`                    // 遊戲類別編號
	Bet03        decimal.Decimal `gorm:"column:bet03;not null;comment:場次編號" json:"bet03"`                      // 場次編號
	Bet04        int             `gorm:"column:bet04;not null;comment:子場次編號" json:"bet04"`                     // 子場次編號
	Bet05        int             `gorm:"column:bet05;not null;comment:會員編號" json:"bet05"`                      // 會員編號
	Bet06        time.Time       `gorm:"column:bet06;not null;comment:開局時間" json:"bet06"`                      // 開局時間
	Bet07        time.Time       `gorm:"column:bet07;not null;comment:帳務日期" json:"bet07"`                      // 帳務日期
	Bet08        time.Time       `gorm:"column:bet08;not null;comment:下注時間" json:"bet08"`                      // 下注時間
	Bet09        string          `gorm:"column:bet09;not null;comment:下注內容" json:"bet09"`                      // 下注內容
	Bet10        int             `gorm:"column:bet10;not null;comment:幣別" json:"bet10"`                        // 幣別
	Bet11        decimal.Decimal `gorm:"column:bet11;not null;comment:匯率" json:"bet11"`                        // 匯率
	Bet12        decimal.Decimal `gorm:"column:bet12;not null;comment:起始點數" json:"bet12"`                      // 起始點數
	Bet12a       decimal.Decimal `gorm:"column:bet12a;not null;comment:起始籌碼" json:"bet12a"`                    // 起始籌碼
	Bet13        decimal.Decimal `gorm:"column:bet13;not null;comment:下注金額" json:"bet13"`                      // 下注金額
	Bet14        decimal.Decimal `gorm:"column:bet14;not null;comment:退水％數" json:"bet14"`                      // 退水％數
	Bet15        int             `gorm:"column:bet15;not null;comment:LV1ID" json:"bet15"`                     // LV1ID
	Bet16        int             `gorm:"column:bet16;not null;comment:LV2ID" json:"bet16"`                     // LV2ID
	Bet17        int             `gorm:"column:bet17;not null;comment:LV3ID" json:"bet17"`                     // LV3ID
	Bet18        int             `gorm:"column:bet18;not null;comment:LV4ID" json:"bet18"`                     // LV4ID
	Bet19        int             `gorm:"column:bet19;not null;comment:LV5ID" json:"bet19"`                     // LV5ID
	Bet20        decimal.Decimal `gorm:"column:bet20;not null;comment:LV1佔成" json:"bet20"`                     // LV1佔成
	Bet21        decimal.Decimal `gorm:"column:bet21;not null;comment:LV2佔成" json:"bet21"`                     // LV2佔成
	Bet22        decimal.Decimal `gorm:"column:bet22;not null;comment:LV3佔成" json:"bet22"`                     // LV3佔成
	Bet23        decimal.Decimal `gorm:"column:bet23;not null;comment:LV4佔成" json:"bet23"`                     // LV4佔成
	Bet24        decimal.Decimal `gorm:"column:bet24;not null;comment:LV5佔成" json:"bet24"`                     // LV5佔成
	Bet25        decimal.Decimal `gorm:"column:bet25;not null;comment:LV1退水%數" json:"bet25"`                   // LV1退水%數
	Bet26        decimal.Decimal `gorm:"column:bet26;not null;comment:LV2退水%數" json:"bet26"`                   // LV2退水%數
	Bet27        decimal.Decimal `gorm:"column:bet27;not null;comment:LV3退水%數" json:"bet27"`                   // LV3退水%數
	Bet28        decimal.Decimal `gorm:"column:bet28;not null;comment:LV4退水%數" json:"bet28"`                   // LV4退水%數
	Bet29        decimal.Decimal `gorm:"column:bet29;not null;comment:LV5退水%數" json:"bet29"`                   // LV5退水%數
	Bet30        string          `gorm:"column:bet30;not null;default:N;comment:取消單" json:"bet30"`             // 取消單
	Bet31        int             `gorm:"column:bet31;not null;comment:桌子編號(Table ID)" json:"bet31"`            // 桌子編號(Table ID)
	Bet32        int             `gorm:"column:bet32;not null;comment:房間編號(Room ID)" json:"bet32"`             // 房間編號(Room ID)
	Betwalletid  string          `gorm:"column:betwalletid;not null;comment:下注單一錢包用" json:"betwalletid"`       // 下注單一錢包用
	Gametype     int             `gorm:"column:gametype;not null;comment:1:網投、2:電投" json:"gametype"`           // 1:網投、2:電投
	Commission   int             `gorm:"column:commission;not null;comment:0.一般,1,免傭" json:"commission"`       // 0.一般,1,免傭
	Category     int             `gorm:"column:category;not null;default:1;comment:1.一般,2.小費" json:"category"` // 1.一般,2.小費
	Eid          int             `gorm:"column:eid;not null;comment:荷官" json:"eid"`                            // 荷官
	Serid        int             `gorm:"column:serid;not null;comment:服務編號" json:"serid"`                      // 服務編號
	IP           string          `gorm:"column:ip;not null;comment:下注IP" json:"ip"`                            // 下注IP
	PartnerBetID string          `gorm:"column:partnerBetId;comment:第三方注單編號" json:"partnerBetId"`              // 第三方注單編號
	GameID       string          `gorm:"column:gameId;comment:第三方遊戲代碼" json:"gameId"`                          // 第三方遊戲代碼
	Updatetime   time.Time       `gorm:"column:updatetime;not null;comment:更新時間" json:"updatetime"`            // 更新時間
}

// TableName Bet01's table name
func (*Bet01) TableName() string {
	return TableNameBet01
}
