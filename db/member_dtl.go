package db

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// MemberDtlDao interface defines CRUD operations for MemberDtl
type MemberDtlDao interface {
	QueryByID(tx *gorm.DB, mid, category int64) (*MemberDtl, error)
	Create(tx *gorm.DB, memberDtl *MemberDtl) error
	DeleteByID(tx *gorm.DB, mid, category int64) error
	Update(tx *gorm.DB, memberDtl *MemberDtl) error
	Updates(tx *gorm.DB, mid, category int64, data map[string]interface{}) error
	QueryByMemberID(tx *gorm.DB, mid int64) ([]*MemberDtl, error)
	CreateMemberDtls(tx *gorm.DB, memberDtls []*MemberDtl) error
}

type memberDtlDao struct{}

func NewMemberDtlDao() MemberDtlDao {
	return &memberDtlDao{}
}

func (dao *memberDtlDao) QueryByID(tx *gorm.DB, mid, category int64) (*MemberDtl, error) {
	ret := MemberDtl{}
	err := tx.Where("mem001 = ? AND mem002 = ?", mid, category).First(&ret).Error
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (dao *memberDtlDao) Create(tx *gorm.DB, memberDtl *MemberDtl) error {
	return tx.Create(memberDtl).Error
}

func (dao *memberDtlDao) DeleteByID(tx *gorm.DB, mid, category int64) error {
	return tx.Where("mem001 = ? AND mem002 = ?", mid, category).Delete(&MemberDtl{}).Error
}

func (dao *memberDtlDao) Update(tx *gorm.DB, memberDtl *MemberDtl) error {
	return tx.Where("mem001 = ? AND mem002 = ?", memberDtl.Mem001, memberDtl.Mem002).Updates(memberDtl).Error
}

func (dao *memberDtlDao) Updates(tx *gorm.DB, mid, category int64, data map[string]interface{}) error {
	return tx.Model(&MemberDtl{}).Where("mem001 = ? AND mem002 = ?", mid, category).Updates(data).Error
}

func (dao *memberDtlDao) QueryByMemberID(tx *gorm.DB, mid int64) ([]*MemberDtl, error) {
	var details []*MemberDtl
	err := tx.Where("mem001 = ?", mid).Find(&details).Error
	if err != nil {
		return nil, err
	}
	return details, nil
}

func (dao *memberDtlDao) CreateMemberDtls(tx *gorm.DB, memberDtls []*MemberDtl) error {
	return tx.Create(memberDtls).Error
}

const TableNameMemberDtl = "member_dtl"

// MemberDtl mapped from table <member_dtl>
type MemberDtl struct {
	Mem001 int64           `gorm:"column:mem001;primaryKey;comment:mid" json:"mem001"`  // mid
	Mem002 int64           `gorm:"column:mem002;primaryKey;comment:類別" json:"mem002"`   // 類別
	Mem003 decimal.Decimal `gorm:"column:mem003;not null;comment:退水" json:"mem003"`     // 退水
	Mem004 int64           `gorm:"column:mem004;not null;comment:LV1ID" json:"mem004"`  // LV1ID
	Mem005 int64           `gorm:"column:mem005;not null;comment:LV2ID" json:"mem005"`  // LV2ID
	Mem006 int64           `gorm:"column:mem006;not null;comment:LV3ID" json:"mem006"`  // LV3ID
	Mem007 int64           `gorm:"column:mem007;not null;comment:LV4ID" json:"mem007"`  // LV4ID
	Mem008 int64           `gorm:"column:mem008;not null;comment:LV5ID" json:"mem008"`  // LV5ID
	Mem009 int64           `gorm:"column:mem009;not null;comment:最大押分05" json:"mem009"` // 最大押分05
	Mem010 int64           `gorm:"column:mem010;not null;comment:最大押分06" json:"mem010"` // 最大押分06
	Mem011 int64           `gorm:"column:mem011;not null" json:"mem011"`
	Mem012 int64           `gorm:"column:mem012;not null;comment:最大可贏金額" json:"mem012"`    // 最大可贏金額
	Mem013 time.Time       `gorm:"column:mem013;not null;comment:可贏可輸起算時間點" json:"mem013"` // 可贏可輸起算時間點
	Mem014 int64           `gorm:"column:mem014;not null;comment:最大可輸金額" json:"mem014"`    // 最大可輸金額
	Mem015 string          `gorm:"column:mem015;not null;comment:新版限額" json:"mem015"`      // 新版限額
	Mem016 decimal.Decimal `gorm:"column:mem016;not null;comment:電投退水" json:"mem016"`      // 電投退水
}

// TableName MemberDtl's table name
func (*MemberDtl) TableName() string {
	return TableNameMemberDtl
}
