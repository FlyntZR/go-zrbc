package db

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Bet02Dao interface {
	QueryByID(tx *gorm.DB, id int64) (*Bet02, error)
	Create(tx *gorm.DB, bet02 *Bet02) (int64, error)
	DeleteByID(tx *gorm.DB, uniqueID int64) error
	Update(tx *gorm.DB, bet02 *Bet02) error
	Updates(tx *gorm.DB, bet01 int64, data map[string]interface{}) error
	GetAgentWinloss(tx *gorm.DB, agentID int64) (decimal.Decimal, error)
}

type bet02Dao struct{}

func NewBet02Dao() Bet02Dao {
	return &bet02Dao{}
}

func (dao *bet02Dao) QueryByID(tx *gorm.DB, id int64) (*Bet02, error) {
	ret := Bet02{}
	err := tx.First(&ret, id).Error
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (dao *bet02Dao) Create(tx *gorm.DB, bet02 *Bet02) (int64, error) {
	err := tx.Table(TableNameBet02).Create(bet02).Error
	if err != nil {
		return 0, err
	}
	return bet02.Bet01, nil
}

func (dao *bet02Dao) DeleteByID(tx *gorm.DB, uniqueID int64) error {
	return tx.Where("bet01 = ?", uniqueID).Delete(&Bet02{}).Error
}

func (dao *bet02Dao) Update(tx *gorm.DB, bet02 *Bet02) error {
	return tx.Table(TableNameBet02).Where("bet01 = ?", bet02.Bet01).Updates(bet02).Error
}

func (dao *bet02Dao) Updates(tx *gorm.DB, bet01 int64, data map[string]interface{}) error {
	return tx.Table(TableNameBet02).Where("bet01 = ?", bet01).Updates(data).Error
}

// GetAgentWinloss gets the current winloss total for an agent
func (dao *bet02Dao) GetAgentWinloss(tx *gorm.DB, agentID int64) (decimal.Decimal, error) {
	var winloss decimal.Decimal
	err := tx.Table(TableNameBet02).Where("bet22 = ?", agentID).Select("SUM(bet17)").Scan(&winloss).Error
	if err != nil {
		return decimal.Zero, err
	}
	return winloss, nil
}

const TableNameBet02 = "bet02"

// Bet02 mapped from table <bet02>
type Bet02 struct {
	Bet01          int64           `gorm:"column:bet01;primaryKey;comment:注單編號" json:"bet01"`                                     // 注單編號
	Bet02          int             `gorm:"column:bet02;not null;comment:遊戲類別編號" json:"bet02"`                                     // 遊戲類別編號
	Bet03          decimal.Decimal `gorm:"column:bet03;not null;comment:場次編號" json:"bet03"`                                       // 場次編號
	Bet04          int             `gorm:"column:bet04;not null;comment:子場次編號" json:"bet04"`                                      // 子場次編號
	Bet05          int             `gorm:"column:bet05;not null;comment:會員編號" json:"bet05"`                                       // 會員編號
	Bet06          time.Time       `gorm:"column:bet06;not null;comment:開局時間" json:"bet06"`                                       // 開局時間
	Bet07          time.Time       `gorm:"column:bet07;not null;comment:帳務日期" json:"bet07"`                                       // 帳務日期
	Bet08          time.Time       `gorm:"column:bet08;not null;comment:下注時間" json:"bet08"`                                       // 下注時間
	Bet09          string          `gorm:"column:bet09;not null;comment:下注內容" json:"bet09"`                                       // 下注內容
	Bet10          int             `gorm:"column:bet10;not null;comment:幣別" json:"bet10"`                                         // 幣別
	Bet11          decimal.Decimal `gorm:"column:bet11;not null;comment:匯率" json:"bet11"`                                         // 匯率
	Bet12          decimal.Decimal `gorm:"column:bet12;not null;comment:起始點數" json:"bet12"`                                       // 起始點數
	Bet12a         decimal.Decimal `gorm:"column:bet12a;not null;comment:起始籌碼" json:"bet12a"`                                     // 起始籌碼
	Bet13          decimal.Decimal `gorm:"column:bet13;not null;comment:下注金額" json:"bet13"`                                       // 下注金額
	Bet14          decimal.Decimal `gorm:"column:bet14;not null;comment:派彩" json:"bet14"`                                         // 派彩
	Bet15          decimal.Decimal `gorm:"column:bet15;not null;comment:退水％數" json:"bet15"`                                       // 退水％數
	Bet16          decimal.Decimal `gorm:"column:bet16;not null;comment:退水金額" json:"bet16"`                                       // 退水金額
	Bet17          decimal.Decimal `gorm:"column:bet17;not null;comment:結果" json:"bet17"`                                         // 結果
	Bet18          int             `gorm:"column:bet18;not null;comment:LV1ID" json:"bet18"`                                      // LV1ID
	Bet19          int             `gorm:"column:bet19;not null;comment:LV2ID" json:"bet19"`                                      // LV2ID
	Bet20          int             `gorm:"column:bet20;not null;comment:LV3ID" json:"bet20"`                                      // LV3ID
	Bet21          int             `gorm:"column:bet21;not null;comment:LV4ID" json:"bet21"`                                      // LV4ID
	Bet22          int             `gorm:"column:bet22;not null;comment:LV5ID" json:"bet22"`                                      // LV5ID
	Bet23          decimal.Decimal `gorm:"column:bet23;not null;comment:LV1佔成" json:"bet23"`                                      // LV1佔成
	Bet24          decimal.Decimal `gorm:"column:bet24;not null;comment:LV2佔成" json:"bet24"`                                      // LV2佔成
	Bet25          decimal.Decimal `gorm:"column:bet25;not null;comment:LV3佔成" json:"bet25"`                                      // LV3佔成
	Bet26          decimal.Decimal `gorm:"column:bet26;not null;comment:LV4佔成" json:"bet26"`                                      // LV4佔成
	Bet27          decimal.Decimal `gorm:"column:bet27;not null;comment:LV5佔成" json:"bet27"`                                      // LV5佔成
	Bet28          decimal.Decimal `gorm:"column:bet28;not null;comment:LV1退水%數" json:"bet28"`                                    // LV1退水%數
	Bet29          decimal.Decimal `gorm:"column:bet29;not null;comment:LV2退水%數" json:"bet29"`                                    // LV2退水%數
	Bet30          decimal.Decimal `gorm:"column:bet30;not null;comment:LV3退水%數" json:"bet30"`                                    // LV3退水%數
	Bet31          decimal.Decimal `gorm:"column:bet31;not null;comment:LV4退水%數" json:"bet31"`                                    // LV4退水%數
	Bet32          decimal.Decimal `gorm:"column:bet32;not null;comment:LV5退水%數" json:"bet32"`                                    // LV5退水%數
	Bet33          decimal.Decimal `gorm:"column:bet33;not null;comment:LV1結果" json:"bet33"`                                      // LV1結果
	Bet34          decimal.Decimal `gorm:"column:bet34;not null;comment:LV2結果" json:"bet34"`                                      // LV2結果
	Bet35          decimal.Decimal `gorm:"column:bet35;not null;comment:LV3結果" json:"bet35"`                                      // LV3結果
	Bet36          decimal.Decimal `gorm:"column:bet36;not null;comment:LV4結果" json:"bet36"`                                      // LV4結果
	Bet37          decimal.Decimal `gorm:"column:bet37;not null;comment:LV5結果" json:"bet37"`                                      // LV5結果
	Bet38          string          `gorm:"column:bet38;not null;default:N;comment:重對" json:"bet38"`                               // 重對
	Bet39          int             `gorm:"column:bet39;not null;comment:桌子編號" json:"bet39"`                                       // 桌子編號
	Bet40          int             `gorm:"column:bet40;not null;comment:房間編號" json:"bet40"`                                       // 房間編號
	Bet41          decimal.Decimal `gorm:"column:bet41;not null;comment:下注退水金額" json:"bet41"`                                     // 下注退水金額
	Betwalletid    string          `gorm:"column:betwalletid;not null;comment:下注單一錢包用" json:"betwalletid"`                        // 下注單一錢包用
	Resultwalletid string          `gorm:"column:resultwalletid;not null;comment:結算單一錢包用" json:"resultwalletid"`                  // 結算單一錢包用
	Validbet       decimal.Decimal `gorm:"column:validbet;not null;comment:有效投注" json:"validbet"`                                 // 有效投注
	Gametype       int             `gorm:"column:gametype;not null;comment:1:網投、2:電投、3 電投網投模式、4 電投電投模式" json:"gametype"`          // 1:網投、2:電投、3 電投網投模式、4 電投電投模式
	Commission     int             `gorm:"column:commission;not null;comment:0.一般,1,免傭" json:"commission"`                        // 0.一般,1,免傭
	Category       int             `gorm:"column:category;not null;default:1;comment:1.一般,2.小費" json:"category"`                  // 1.一般,2.小費
	Eid            int             `gorm:"column:eid;not null;comment:荷官" json:"eid"`                                             // 荷官
	Serid          int             `gorm:"column:serid;not null;comment:服務編號" json:"serid"`                                       // 服務編號
	IP             string          `gorm:"column:ip;not null;comment:下注IP" json:"ip"`                                             // 下注IP
	PartnerBetID   string          `gorm:"column:partnerBetId;comment:第三方注單編號" json:"partnerBetId"`                               // 第三方注單編號
	GameID         string          `gorm:"column:gameId;comment:第三方遊戲下注代號或名稱" json:"gameId"`                                      // 第三方遊戲下注代號或名稱
	Updatetime     time.Time       `gorm:"column:updatetime;not null;default:current_timestamp();comment:更新時間" json:"updatetime"` // 更新時間
}

// TableName Bet02's table name
func (*Bet02) TableName() string {
	return TableNameBet02
}
