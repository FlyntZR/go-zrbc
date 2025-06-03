package db

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type UserDao interface {
	QueryByID(tx *gorm.DB, id int64) (*Member, error)
	QueryByAccountAndPwd(tx *gorm.DB, account, passwd string) (*Member, error)
	CreateUser(tx *gorm.DB, member *Member) (int64, error)
	DeleteByID(tx *gorm.DB, uniqueID int64) error
	UpdateMember(tx *gorm.DB, ysUser *Member) error
	UpdatesMember(tx *gorm.DB, userID int64, data map[string]interface{}) error
}

type userDao struct{}

func NewMemberDao() UserDao {
	return &userDao{}
}

func (dao *userDao) QueryByID(tx *gorm.DB, id int64) (*Member, error) {
	ret := Member{}
	err := tx.First(&ret, id).Error
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (dao *userDao) QueryByAccountAndPwd(tx *gorm.DB, account, passwd string) (*Member, error) {
	ret := Member{}
	err := tx.Where("mem002 = ? AND mem003 = ?", account, passwd).First(&ret).Error
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (dao *userDao) CreateUser(tx *gorm.DB, member *Member) (int64, error) {
	err := tx.Table("member").Create(member).Error
	if err != nil {
		return 0, err
	}

	return member.ID, nil
}

func (dao *userDao) DeleteByID(tx *gorm.DB, uniqueID int64) error {
	return tx.Where("mem001 = ?", uniqueID).Delete(&Member{}).Error
}

func (dao *userDao) UpdateMember(tx *gorm.DB, ysUser *Member) error {
	return tx.Table("member").Where("mem001 = ?", ysUser.ID).Updates(ysUser).Error
}

func (dao *userDao) UpdatesMember(tx *gorm.DB, userID int64, data map[string]interface{}) error {
	return tx.Table("member").Where("mem001 = ?", userID).Updates(data).Error
}

const TableNameMember = "member"

// Member mapped from table <member>
type Member struct {
	ID              int64           `gorm:"column:mem001;primaryKey;autoIncrement:true" json:"mem001"` // 主键
	User            string          `gorm:"column:mem002;not null" json:"mem002"`                      // 用户名
	Password        string          `gorm:"column:mem003;not null" json:"mem003"`                      // 用户密码
	UserName        string          `gorm:"column:mem004;not null" json:"mem004"`                      // 用户名称
	Mem005          time.Time       `gorm:"column:mem005;not null;default:current_timestamp()" json:"mem005"`
	Mem006          int             `gorm:"column:mem006;not null" json:"mem006"`
	Mem007          int             `gorm:"column:mem007;not null" json:"mem007"`
	Mem008          int             `gorm:"column:mem008;not null" json:"mem008"`
	Mem009          int             `gorm:"column:mem009;not null" json:"mem009"`
	Mem010          int             `gorm:"column:mem010;not null" json:"mem010"`
	Mem011          int             `gorm:"column:mem011;not null" json:"mem011"`
	Mem012          int             `gorm:"column:mem012;not null" json:"mem012"`
	Mem013          time.Time       `gorm:"column:mem013;not null" json:"mem013"`
	Mem014          string          `gorm:"column:mem014;not null" json:"mem014"`
	Mem015          int             `gorm:"column:mem015;not null;comment:login_error" json:"mem015"`         // login_error
	Mem016          string          `gorm:"column:mem016;not null;default:Y;comment:enable" json:"mem016"`    // enable
	Mem017          string          `gorm:"column:mem017;not null;default:Y;comment:canbet" json:"mem017"`    // canbet
	Mem018          string          `gorm:"column:mem018;not null;default:Y;comment:chg_pw" json:"mem018"`    // chg_pw
	Mem019          string          `gorm:"column:mem019;not null;default:N;comment:is_test" json:"mem019"`   // is_test
	Mem020          string          `gorm:"column:mem020;not null;default:N;comment:be_traded" json:"mem020"` // be_traded
	Mem021          int             `gorm:"column:mem021;not null" json:"mem021"`
	Mem022          string          `gorm:"column:mem022;not null;comment:電話" json:"mem022"`     // 電話
	Mem022a         int             `gorm:"column:mem022a;not null;comment:電話簡碼" json:"mem022a"` // 電話簡碼
	Mem023          int             `gorm:"column:mem023;not null" json:"mem023"`
	Mem024          int             `gorm:"column:mem024;not null;comment:mem_risk" json:"mem024"` // mem_risk
	Mem026          int             `gorm:"column:mem026;not null" json:"mem026"`
	Mem028          string          `gorm:"column:mem028;not null;comment:備註" json:"mem028"`         // 備註
	Type            int             `gorm:"column:type;not null;comment:0:現金;1:信用;2:電投" json:"type"` // 0:現金;1:信用;2:電投
	Currency        int             `gorm:"column:currency;not null;comment:幣別" json:"currency"`     // 幣別
	Cash            decimal.Decimal `gorm:"column:cash;not null;comment:現金" json:"cash"`             // 現金
	Money           decimal.Decimal `gorm:"column:money;not null;comment:己出碼額度" json:"money"`        // 己出碼額度
	Lockmoney       decimal.Decimal `gorm:"column:lockmoney;not null;comment:鎖定金額" json:"lockmoney"` // 鎖定金額
	Head            int             `gorm:"column:head;not null;comment:頭像ID" json:"head"`           // 頭像ID
	Chips           string          `gorm:"column:chips;not null;comment:籌碼選擇" json:"chips"`         // 籌碼選擇
	Follow1         string          `gorm:"column:follow1;not null;comment:注關會員" json:"follow1"`     // 注關會員
	Follow2         string          `gorm:"column:follow2;not null;comment:關注荷官" json:"follow2"`     // 關注荷官
	Tip             string          `gorm:"column:tip;not null;default:Y" json:"tip"`
	Red             string          `gorm:"column:red;not null;default:N;comment:红包开关" json:"red"` // 红包开关
	Wallet          string          `gorm:"column:wallet;not null;default:N" json:"wallet"`
	Opengame        string          `gorm:"column:opengame;not null;comment:輸入要開啟的種類" json:"opengame"` // 輸入要開啟的種類
	Site            string          `gorm:"column:site;not null" json:"site"`
	Lineid          string          `gorm:"column:lineid;not null;comment:Line ID" json:"lineid"`                            // Line ID
	Kickperiod      int             `gorm:"column:kickperiod;not null;default:10;comment:踢出局數" json:"kickperiod"`            // 踢出局數
	Identity        int             `gorm:"column:identity;not null;comment:0:信用,1:api" json:"identity"`                     // 0:信用,1:api
	Singlebetprompt decimal.Decimal `gorm:"column:singlebetprompt;not null;default:0;comment:單注金額告警" json:"singlebetprompt"` // 單注金額告警
	Conwinprompt    int             `gorm:"column:conwinprompt;not null;comment:連贏局數告警" json:"conwinprompt"`                 // 連贏局數告警
	Winlossprompt   string          `gorm:"column:winlossprompt;not null;default:0;comment:最大可輸可贏提示" json:"winlossprompt"`   // 最大可輸可贏提示
	Onlineprompt    string          `gorm:"column:onlineprompt;not null;comment:上線提示" json:"onlineprompt"`                   // 上線提示
	Profitprompt    int             `gorm:"column:profitprompt;not null;comment:盈利告警" json:"profitprompt"`                   // 盈利告警
}

// TableName Member's table name
func (*Member) TableName() string {
	return TableNameMember
}
