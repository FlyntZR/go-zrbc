package db

import (
	"time"

	"gorm.io/gorm"
)

// MemLoginDao defines the interface for MemLogin database operations
type MemLoginDao interface {
	QueryByID(tx *gorm.DB, memID int64) (*MemLogin, error)
	QueryBySID(tx *gorm.DB, sid string) (*MemLogin, error)
	Create(tx *gorm.DB, memLogin *MemLogin) error
	DeleteByID(tx *gorm.DB, memID int64) error
	Update(tx *gorm.DB, memLogin *MemLogin) error
	UpdateFields(tx *gorm.DB, memID int64, data map[string]interface{}) error
	CreateOrUpdateMemLogin(tx *gorm.DB, uid int64, wcode int, sid string, userIP string, now time.Time) (*MemLogin, error)
}

// memLoginDao implements MemLoginDao interface
type memLoginDao struct{}

// NewMemLoginDao creates a new instance of MemLoginDao
func NewMemLoginDao() MemLoginDao {
	return &memLoginDao{}
}

func (dao *memLoginDao) QueryByID(tx *gorm.DB, memID int64) (*MemLogin, error) {
	ret := MemLogin{}
	err := tx.Where("mlg001 = ?", memID).First(&ret).Error
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (dao *memLoginDao) QueryBySID(tx *gorm.DB, sid string) (*MemLogin, error) {
	ret := MemLogin{}
	err := tx.Where("mlg003 = ?", sid).First(&ret).Error
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (dao *memLoginDao) Create(tx *gorm.DB, memLogin *MemLogin) error {
	return tx.Create(memLogin).Error
}

func (dao *memLoginDao) DeleteByID(tx *gorm.DB, memID int64) error {
	return tx.Where("mlg001 = ?", memID).Delete(&MemLogin{}).Error
}

func (dao *memLoginDao) Update(tx *gorm.DB, memLogin *MemLogin) error {
	return tx.Where("mlg001 = ?", memLogin.Mlg001).Updates(memLogin).Error
}

func (dao *memLoginDao) UpdateFields(tx *gorm.DB, memID int64, data map[string]interface{}) error {
	return tx.Model(&MemLogin{}).Where("mlg001 = ?", memID).Updates(data).Error
}

// ProWriteMemLogin handles member login record creation/update and updates member's last login info
func (dao *memLoginDao) CreateOrUpdateMemLogin(tx *gorm.DB, uid int64, wcode int, sid string, userIP string, now time.Time) (*MemLogin, error) {
	// First check if a record exists
	var existingLogin MemLogin
	result := tx.Where("mlg001 = ? AND mlg002 = ?", uid, wcode).First(&existingLogin)

	if result.Error == gorm.ErrRecordNotFound {
		// Insert new record
		newLogin := &MemLogin{
			Mlg001: uid,
			Mlg002: wcode,
			Mlg003: sid,
			Mlg004: now,
			Mlg005: now, // For new records, last login is same as this login
			Mlg006: userIP,
		}
		if err := tx.Create(newLogin).Error; err != nil {
			return nil, err
		}
		return newLogin, nil
	} else if result.Error == nil {
		// Update existing record
		updates := map[string]interface{}{
			"mlg003": sid,
			"mlg005": existingLogin.Mlg004, // Set last login to previous this_login
			"mlg004": now,
			"mlg006": userIP,
		}
		if err := tx.Model(&MemLogin{}).Where("mlg001 = ? AND mlg002 = ?", uid, wcode).Updates(updates).Error; err != nil {
			return nil, err
		}
		existingLogin.Mlg003 = sid
		existingLogin.Mlg005 = existingLogin.Mlg004
		existingLogin.Mlg004 = now
		existingLogin.Mlg006 = userIP
		return &existingLogin, nil
	} else {
		return nil, result.Error
	}
}

const TableNameMemLogin = "mem_login"

// MemLogin mapped from table <mem_login>
type MemLogin struct {
	Mlg001 int64     `gorm:"column:mlg001;primaryKey;comment:mem_id" json:"mlg001"`                               // mem_id
	Mlg002 int       `gorm:"column:mlg002;not null;comment:site" json:"mlg002"`                                   // site
	Mlg003 string    `gorm:"column:mlg003;not null;comment:sid" json:"mlg003"`                                    // sid
	Mlg004 time.Time `gorm:"column:mlg004;not null;default:current_timestamp();comment:this_login" json:"mlg004"` // this_login
	Mlg005 time.Time `gorm:"column:mlg005;not null;comment:last_login" json:"mlg005"`                             // last_login
	Mlg006 string    `gorm:"column:mlg006;not null;comment:ip" json:"mlg006"`                                     // ip
	Mlg007 string    `gorm:"column:mlg007;not null;comment:info" json:"mlg007"`                                   // info
	Mlg008 int       `gorm:"column:mlg008;not null;comment:on_line" json:"mlg008"`                                // on_line
}

// TableName MemLogin's table name
func (*MemLogin) TableName() string {
	return TableNameMemLogin
}
