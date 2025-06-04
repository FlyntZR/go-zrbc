package db

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type AgentDao interface {
	QueryByID(tx *gorm.DB, id int64) (*Agent, error)
	QueryByVendorID(tx *gorm.DB, vendorID string) (*Agent, error)
	CreateAgent(tx *gorm.DB, agent *Agent) (int64, error)
	DeleteByID(tx *gorm.DB, uniqueID int64) error
	UpdateAgent(tx *gorm.DB, agent *Agent) error
	UpdatesAgent(tx *gorm.DB, agentID int64, data map[string]interface{}) error
}

type agentDao struct{}

func NewAgentDao() AgentDao {
	return &agentDao{}
}

func (dao *agentDao) QueryByID(tx *gorm.DB, id int64) (*Agent, error) {
	ret := Agent{}
	err := tx.First(&ret, id).Error
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (dao *agentDao) QueryByVendorID(tx *gorm.DB, vendorID string) (*Agent, error) {
	ret := Agent{}
	err := tx.Where("age002 = ?", vendorID).First(&ret).Error
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (dao *agentDao) CreateAgent(tx *gorm.DB, agent *Agent) (int64, error) {
	err := tx.Table(TableNameAgent).Create(agent).Error
	if err != nil {
		return 0, err
	}
	return agent.Age001, nil
}

func (dao *agentDao) DeleteByID(tx *gorm.DB, uniqueID int64) error {
	return tx.Where("age001 = ?", uniqueID).Delete(&Agent{}).Error
}

func (dao *agentDao) UpdateAgent(tx *gorm.DB, agent *Agent) error {
	return tx.Table(TableNameAgent).Where("age001 = ?", agent.Age001).Updates(agent).Error
}

func (dao *agentDao) UpdatesAgent(tx *gorm.DB, agentID int64, data map[string]interface{}) error {
	return tx.Table(TableNameAgent).Where("age001 = ?", agentID).Updates(data).Error
}

const TableNameAgent = "agent"

// Agent mapped from table <agent>
type Agent struct {
	Age001       int64           `gorm:"column:age001;primaryKey;autoIncrement:true;comment:sn" json:"age001"` // sn
	Age002       string          `gorm:"column:age002;not null" json:"age002"`
	Age003       string          `gorm:"column:age003;not null" json:"age003"`
	Age004       string          `gorm:"column:age004;not null" json:"age004"`
	Age005       time.Time       `gorm:"column:age005;not null;default:current_timestamp()" json:"age005"`
	Age006       int             `gorm:"column:age006;not null" json:"age006"`
	Age007       int             `gorm:"column:age007;not null" json:"age007"`
	Age008       int             `gorm:"column:age008;not null" json:"age008"`
	Age009       int             `gorm:"column:age009;not null" json:"age009"`
	Age010       int             `gorm:"column:age010;not null" json:"age010"`
	Age011       int             `gorm:"column:age011;not null" json:"age011"`
	Age012       time.Time       `gorm:"column:age012;not null" json:"age012"`
	Age013       string          `gorm:"column:age013;not null" json:"age013"`
	Age014       int             `gorm:"column:age014;not null;comment:風險重置開關" json:"age014"`          // 風險重置開關
	Age015       string          `gorm:"column:age015;not null;default:Y;comment:啟停用狀態" json:"age015"` // 啟停用狀態
	Age016       string          `gorm:"column:age016;not null;default:Y;comment:啟停押狀態" json:"age016"` // 啟停押狀態
	Age017       string          `gorm:"column:age017;not null;default:Y" json:"age017"`
	Age018       string          `gorm:"column:age018;not null;default:N" json:"age018"`
	Age019       int             `gorm:"column:age019;not null;default:1;comment:結算報表開關" json:"age019"` // 結算報表開關
	Age020       int             `gorm:"column:age020;not null;comment:結算報表格式" json:"age020"`           // 結算報表格式
	Age021       int             `gorm:"column:age021;not null;comment:結算報表語系" json:"age021"`           // 結算報表語系
	Age022       int             `gorm:"column:age022;not null;comment:測試線判別" json:"age022"`            // 測試線判別
	Age023       string          `gorm:"column:age023;not null" json:"age023"`
	Age024       string          `gorm:"column:age024;not null;comment:備註" json:"age024"` // 備註
	Credit       decimal.Decimal `gorm:"column:credit;not null" json:"credit"`
	Cash         decimal.Decimal `gorm:"column:cash;not null" json:"cash"`
	Type         int             `gorm:"column:type;not null" json:"type"`
	Currency     int             `gorm:"column:currency;not null;default:1" json:"currency"`
	Tip          string          `gorm:"column:tip;not null;default:Y" json:"tip"`
	Red          string          `gorm:"column:red;not null;default:N;comment:紅包开关" json:"red"` // 紅包开关
	PrefixAdd    string          `gorm:"column:prefix_add;not null;default:N" json:"prefix_add"`
	PrefixAcc    string          `gorm:"column:prefix_acc;not null;comment:前綴碼" json:"prefix_acc"` // 前綴碼
	ReceiptAcc   int             `gorm:"column:receipt_acc;comment:收款帳號" json:"receipt_acc"`       // 收款帳號
	Notification int             `gorm:"column:notification;comment:通知帳號" json:"notification"`     // 通知帳號
	Sacc         string          `gorm:"column:sacc;not null;comment:顯示帳號" json:"sacc"`            // 顯示帳號
	Opengame     string          `gorm:"column:opengame;not null" json:"opengame"`
	Site         string          `gorm:"column:site;not null" json:"site"`
	Membermax    int             `gorm:"column:membermax;not null;comment:會員總數" json:"membermax"`              // 會員總數
	Profitmax    string          `gorm:"column:profitmax;not null;default:0;comment:最大利潤" json:"profitmax"`    // 最大利潤
	Promotecode  string          `gorm:"column:promotecode;not null;comment:推廣代碼" json:"promotecode"`          // 推廣代碼
	Kickperiod   int             `gorm:"column:kickperiod;not null;default:10;comment:踢出局數" json:"kickperiod"` // 踢出局數
	Identity     int             `gorm:"column:identity;not null;comment:0:信用;1:API" json:"identity"`          // 0:信用;1:API
	ChkKey       string          `gorm:"column:chkKey;not null" json:"chkKey"`
	ChkLock      string          `gorm:"column:chkLock;not null;default:N;comment:是否鎖定" json:"chkLock"` // 是否鎖定
	LastChpsw    time.Time       `gorm:"column:last_chpsw;not null;comment:上次修改密碼時間" json:"last_chpsw"` // 上次修改密碼時間
}

// TableName Agent's table name
func (*Agent) TableName() string {
	return TableNameAgent
}
