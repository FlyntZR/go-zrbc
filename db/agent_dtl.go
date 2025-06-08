package db

import (
	"fmt"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type AgentDtlDao interface {
	QueryByAgentID(tx *gorm.DB, agentID int64) ([]*AgentDtl, error)
	QueryByAgentIDAndCategory(tx *gorm.DB, agentID int64, category int64) (*AgentDtl, error)
	Create(tx *gorm.DB, agentDtl *AgentDtl) error
	DeleteByID(tx *gorm.DB, id int64, category int64) error
	Update(tx *gorm.DB, agentDtl *AgentDtl) error
	Updates(tx *gorm.DB, id int64, category int64, data map[string]interface{}) error
	GetDefaultLimits(tx *gorm.DB, agentID int64) (map[string]string, error)
}

type agentDtlDao struct{}

func NewAgentDtlDao() AgentDtlDao {
	return &agentDtlDao{}
}

func (dao *agentDtlDao) QueryByAgentID(tx *gorm.DB, agentID int64) ([]*AgentDtl, error) {
	ret := []*AgentDtl{}
	err := tx.Where("ag001 = ?", agentID).Find(&ret).Error
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (dao *agentDtlDao) QueryByAgentIDAndCategory(tx *gorm.DB, agentID int64, category int64) (*AgentDtl, error) {
	ret := AgentDtl{}
	err := tx.Where("ag001 = ? AND ag002 = ?", agentID, category).First(&ret).Error
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (dao *agentDtlDao) Create(tx *gorm.DB, agentDtl *AgentDtl) error {
	return tx.Table(TableNameAgentDtl).Create(agentDtl).Error
}

func (dao *agentDtlDao) DeleteByID(tx *gorm.DB, id int64, category int64) error {
	return tx.Where("ag001 = ? AND ag002 = ?", id, category).Delete(&AgentDtl{}).Error
}

func (dao *agentDtlDao) Update(tx *gorm.DB, agentDtl *AgentDtl) error {
	return tx.Table(TableNameAgentDtl).Where("ag001 = ? AND ag002 = ?", agentDtl.Ag001, agentDtl.Ag002).Updates(agentDtl).Error
}

func (dao *agentDtlDao) Updates(tx *gorm.DB, id int64, category int64, data map[string]interface{}) error {
	return tx.Table(TableNameAgentDtl).Where("ag001 = ? AND ag002 = ?", id, category).Updates(data).Error
}

// GetDefaultLimits gets the default game limits for an agent
func (dao *agentDtlDao) GetDefaultLimits(tx *gorm.DB, agentID int64) (map[string]string, error) {
	agentDtlList := []*AgentDtl{}
	err := tx.Table(TableNameAgentDtl).Where("ag001 = ? AND ag002 IN('101','102','103','104','105','106','107','108','110','111','112','113','117','121','126')", agentID).Find(&agentDtlList).Error
	if err != nil {
		return nil, err
	}
	limits := make(map[string]string)
	for _, agentDtl := range agentDtlList {
		limits[fmt.Sprintf("set%d_14", agentDtl.Ag002)] = agentDtl.Ag014
	}

	return limits, nil
}

const TableNameAgentDtl = "agent_dtl"

// AgentDtl mapped from table <agent_dtl>
type AgentDtl struct {
	Ag001 int64           `gorm:"column:ag001;primaryKey;autoIncrement:true;comment:id" json:"ag001"` // id
	Ag002 int64           `gorm:"column:ag002;primaryKey;comment:類別" json:"ag002"`                    // 類別
	Ag003 decimal.Decimal `gorm:"column:ag003;not null;comment:退水" json:"ag003"`                      // 退水
	Ag004 int             `gorm:"column:ag004;not null;comment:最小押分" json:"ag004"`                    // 最小押分
	Ag005 int             `gorm:"column:ag005;not null;comment:最大押分" json:"ag005"`                    // 最大押分
	Ag006 int             `gorm:"column:ag006;not null;comment:押分倍數" json:"ag006"`                    // 押分倍數
	Ag007 int             `gorm:"column:ag007;not null;comment:最大押分02" json:"ag007"`                  // 最大押分02
	Ag008 int             `gorm:"column:ag008;not null;comment:最大押分03" json:"ag008"`                  // 最大押分03
	Ag009 int             `gorm:"column:ag009;not null;comment:最大押分04" json:"ag009"`                  // 最大押分04
	Ag010 int             `gorm:"column:ag010;not null;comment:最大押分05" json:"ag010"`                  // 最大押分05
	Ag011 int             `gorm:"column:ag011;not null;comment:最大押分06" json:"ag011"`                  // 最大押分06
	Ag012 decimal.Decimal `gorm:"column:ag012;not null;comment:佔成" json:"ag012"`                      // 佔成
	Ag013 int             `gorm:"column:ag013;not null;comment:選則限紅（BIT）" json:"ag013"`               // 選則限紅（BIT）
	Ag014 string          `gorm:"column:ag014;not null;comment:新版限額" json:"ag014"`                    // 新版限額
	Ag015 decimal.Decimal `gorm:"column:ag015;not null;comment:電投退水" json:"ag015"`                    // 電投退水
	Ag016 int             `gorm:"column:ag016;not null;comment:電投佔成" json:"ag016"`                    // 電投佔成
}

// TableName AgentDtl's table name
func (*AgentDtl) TableName() string {
	return TableNameAgentDtl
}
