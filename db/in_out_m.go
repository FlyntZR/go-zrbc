package db

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type InOutMDao interface {
	QueryByID(tx *gorm.DB, id int64) (*InOutM, error)
	QueryByMemberID(tx *gorm.DB, memberID int64) ([]*InOutM, error)
	Create(tx *gorm.DB, record *InOutM) (int64, error)
	DeleteByID(tx *gorm.DB, id int64) error
	Update(tx *gorm.DB, record *InOutM) error
	Updates(tx *gorm.DB, id int64, data map[string]interface{}) error
	CountTransactionsInLastMinute(tx *gorm.DB, memberID int64) (int64, error)
	GetLastTransaction(tx *gorm.DB, memberID int64, orderNum string) (*InOutM, error)
}

type inOutMDao struct{}

func NewInOutMDao() InOutMDao {
	return &inOutMDao{}
}

// LastTransactionResult represents the result of the last transaction query
type LastTransactionResult struct {
	OrderID  int64     `gorm:"column:orderid"`
	LastTime time.Time `gorm:"column:lasttime"`
	Money    string    `gorm:"column:money"`
	OrderNum string    `gorm:"column:ordernum"`
}

// GetLastTransaction gets the last transaction for a member, optionally filtered by order number
func (dao *inOutMDao) GetLastTransaction(tx *gorm.DB, memberID int64, orderNum string) (*InOutM, error) {
	query := tx.Table(TableNameInOutM).Where("iom003 = ?", memberID)
	if orderNum != "" {
		query = query.Where("iom008 = ?", orderNum)
	}
	var result InOutM
	err := query.Order("iom002 DESC").First(&result).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (dao *inOutMDao) QueryByID(tx *gorm.DB, id int64) (*InOutM, error) {
	ret := InOutM{}
	err := tx.First(&ret, id).Error
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (dao *inOutMDao) QueryByMemberID(tx *gorm.DB, memberID int64) ([]*InOutM, error) {
	var records []*InOutM
	err := tx.Where("iom003 = ?", memberID).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (dao *inOutMDao) Create(tx *gorm.DB, record *InOutM) (int64, error) {
	err := tx.Table(TableNameInOutM).Create(record).Error
	if err != nil {
		return 0, err
	}
	return record.Iom001, nil
}

func (dao *inOutMDao) DeleteByID(tx *gorm.DB, id int64) error {
	return tx.Where("iom001 = ?", id).Delete(&InOutM{}).Error
}

func (dao *inOutMDao) Update(tx *gorm.DB, record *InOutM) error {
	return tx.Table(TableNameInOutM).Where("iom001 = ?", record.Iom001).Updates(record).Error
}

func (dao *inOutMDao) Updates(tx *gorm.DB, id int64, data map[string]interface{}) error {
	return tx.Table(TableNameInOutM).Where("iom001 = ?", id).Updates(data).Error
}

// CountTransactionsInLastMinute counts the number of transactions for a member in the last minute
func (dao *inOutMDao) CountTransactionsInLastMinute(tx *gorm.DB, memberID int64) (int64, error) {
	var count int64
	lastTime := time.Now().Add(-1 * time.Minute)
	nowTime := time.Now()

	err := tx.Raw(`
		SELECT count(1) 
		FROM in_out_m 
		WHERE iom003 = ? 
		AND iom002 BETWEEN ? AND ? 
		AND iom005 in (121,122)`,
		memberID,
		lastTime.Format("2006-01-02 15:04:05"),
		nowTime.Format("2006-01-02 15:04:05"),
	).Scan(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

const TableNameInOutM = "in_out_m"

// InOutM mapped from table <in_out_m>
type InOutM struct {
	Iom001 int64           `gorm:"column:iom001;primaryKey;autoIncrement:true;comment:sn" json:"iom001"`          // sn
	Iom002 time.Time       `gorm:"column:iom002;not null;default:current_timestamp();comment:time" json:"iom002"` // time
	Iom003 int64           `gorm:"column:iom003;not null;comment:member" json:"iom003"`                           // member
	Iom004 decimal.Decimal `gorm:"column:iom004;not null;comment:money" json:"iom004"`                            // money
	Iom005 string          `gorm:"column:iom005;not null;comment:op_code" json:"iom005"`                          // op_code
	Iom006 int64           `gorm:"column:iom006;not null;comment:op_lv" json:"iom006"`                            // op_lv
	Iom007 int64           `gorm:"column:iom007;not null;comment:op_aid" json:"iom007"`                           // op_aid
	Iom008 string          `gorm:"column:iom008;not null;comment:memo" json:"iom008"`                             // memo
	Iom009 int64           `gorm:"column:iom009;not null;comment:site" json:"iom009"`                             // site
	Iom010 decimal.Decimal `gorm:"column:iom010;not null;comment:subtotal" json:"iom010"`                         // subtotal
}

// TableName InOutM's table name
func (*InOutM) TableName() string {
	return TableNameInOutM
}
