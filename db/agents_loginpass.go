package db

import (
	"time"

	"gorm.io/gorm"
)

type AgentsLoginPassDao interface {
	QueryByID(tx *gorm.DB, id int64) (*AgentsLoginPass, error)
	QueryByAidAndVendorID(tx *gorm.DB, aid int64, vendorID string) (*AgentsLoginPass, error)
	Create(tx *gorm.DB, agent *AgentsLoginPass) (int64, error)
	DeleteByID(tx *gorm.DB, id int64) error
	Update(tx *gorm.DB, agent *AgentsLoginPass) error
	UpdateFields(tx *gorm.DB, id int64, data map[string]interface{}) error
	GetLoginPassByVendor(tx *gorm.DB, vendorID string) (map[string]string, error)
}

type agentsLoginPassDao struct{}

func NewAgentsLoginPassDao() AgentsLoginPassDao {
	return &agentsLoginPassDao{}
}

func (dao *agentsLoginPassDao) QueryByID(tx *gorm.DB, id int64) (*AgentsLoginPass, error) {
	ret := AgentsLoginPass{}
	err := tx.First(&ret, id).Error
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (dao *agentsLoginPassDao) QueryByAidAndVendorID(tx *gorm.DB, aid int64, vendorID string) (*AgentsLoginPass, error) {
	ret := AgentsLoginPass{}
	err := tx.Where("aid = ? AND vendorId = ?", aid, vendorID).First(&ret).Error
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (dao *agentsLoginPassDao) Create(tx *gorm.DB, agent *AgentsLoginPass) (int64, error) {
	err := tx.Create(agent).Error
	if err != nil {
		return 0, err
	}
	return agent.ID, nil
}

func (dao *agentsLoginPassDao) DeleteByID(tx *gorm.DB, id int64) error {
	return tx.Where("id = ?", id).Delete(&AgentsLoginPass{}).Error
}

func (dao *agentsLoginPassDao) Update(tx *gorm.DB, agent *AgentsLoginPass) error {
	return tx.Where("id = ?", agent.ID).Updates(agent).Error
}

func (dao *agentsLoginPassDao) UpdateFields(tx *gorm.DB, id int64, data map[string]interface{}) error {
	return tx.Model(&AgentsLoginPass{}).Where("id = ?", id).Updates(data).Error
}

// GetLoginPassByVendor retrieves login pass information for a vendor
func (dao *agentsLoginPassDao) GetLoginPassByVendor(tx *gorm.DB, vendorID string) (map[string]string, error) {
	var agent AgentsLoginPass
	err := tx.Where("vendorId = ?", vendorID).Select("co").First(&agent).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return map[string]string{}, nil
		}
		return nil, err
	}

	return map[string]string{"co": agent.Co}, nil
}

const TableNameAgentsLoginPass = "agents_LoginPass"

// AgentsLoginPass mapped from table <agents_LoginPass>
type AgentsLoginPass struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Aid          int64     `gorm:"column:aid;not null;comment:agent_id" json:"aid"`           // agent_id
	VendorID     string    `gorm:"column:vendorId;not null;comment:lv5 代理商" json:"vendorId"`  // lv5 代理商
	Signature    string    `gorm:"column:signature;not null;comment:密鑰" json:"signature"`     // 密鑰
	Signature2   string    `gorm:"column:signature2;not null;comment:客戶密鑰" json:"signature2"` // 客戶密鑰
	Password     string    `gorm:"column:password;not null" json:"password"`
	Addtime      time.Time `gorm:"column:addtime;not null" json:"addtime"`
	URL          string    `gorm:"column:url;not null;comment:單一錢包回傳網址" json:"url"` // 單一錢包回傳網址
	Skyname      string    `gorm:"column:skyname;not null" json:"skyname"`
	Type         string    `gorm:"column:type;not null;default:c;comment:c:一般,w:單一" json:"type"`        // c:一般,w:單一
	Lang         int       `gorm:"column:lang;not null;comment:0:中文 ,1:英文" json:"lang"`                 // 0:中文 ,1:英文
	Betfeedback  int       `gorm:"column:betfeedback;not null;comment:異常鎖定,0:N;1:Y" json:"betfeedback"` // 異常鎖定,0:N;1:Y
	GatewayURL   string    `gorm:"column:Gateway_url;not null;comment:正機白名單使用的url" json:"Gateway_url"`  // 正機白名單使用的url
	WhiteList    string    `gorm:"column:whiteList;not null;comment:白名單" json:"whiteList"`              // 白名單
	Operator     string    `gorm:"column:operator;not null" json:"operator"`
	ModifyTime   time.Time `gorm:"column:modify_time;not null;default:current_timestamp()" json:"modify_time"`
	Object       int       `gorm:"column:object;not null;default:1;comment:0:呼叫php,1:呼叫客戶" json:"object"` // 0:呼叫php,1:呼叫客戶
	Settle       int       `gorm:"column:settle;not null;default:1;comment:單一錢包結算" json:"settle"`         // 單一錢包結算
	Co           string    `gorm:"column:co;not null" json:"co"`
	PrefixSwitch string    `gorm:"column:prefix_switch;not null;default:N;comment:前綴碼開關" json:"prefix_switch"` // 前綴碼開關
	OpenGameURL  string    `gorm:"column:openGame_url;not null;comment:指定網址" json:"openGame_url"`              // 指定網址
	Subdomain    string    `gorm:"column:subdomain;not null" json:"subdomain"`
}

// TableName AgentsLoginPass's table name
func (*AgentsLoginPass) TableName() string {
	return TableNameAgentsLoginPass
}
