package db

import (
	"gorm.io/gorm"
)

const TableNameApiurl = "apiurl"

// Apiurl mapped from table <apiurl>
type Apiurl struct {
	ID   int64  `gorm:"column:id;primaryKey;autoIncrement:true;comment:序号" json:"id"` // 序号
	Code int    `gorm:"column:code;not null;comment:参数" json:"code"`                  // 参数
	URL  string `gorm:"column:url;not null;comment:网址" json:"url"`                    // 网址
}

// TableName Apiurl's table name
func (*Apiurl) TableName() string {
	return TableNameApiurl
}

// ApiurlDao defines the interface for Apiurl database operations
type ApiurlDao interface {
	QueryByID(tx *gorm.DB, id int64) (*Apiurl, error)
	QueryByCode(tx *gorm.DB, code int) (*Apiurl, error)
	Create(tx *gorm.DB, apiurl *Apiurl) (int64, error)
	DeleteByID(tx *gorm.DB, id int64) error
	Update(tx *gorm.DB, apiurl *Apiurl) error
	UpdateFields(tx *gorm.DB, id int64, data map[string]interface{}) error
}

// apiurlDaoImpl implements ApiurlDao interface
type apiurlDaoImpl struct {
}

// NewApiurlDao creates a new instance of ApiurlDao
func NewApiurlDao() ApiurlDao {
	return &apiurlDaoImpl{}
}

func (dao *apiurlDaoImpl) QueryByID(tx *gorm.DB, id int64) (*Apiurl, error) {
	ret := Apiurl{}
	err := tx.First(&ret, id).Error
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (dao *apiurlDaoImpl) QueryByCode(tx *gorm.DB, code int) (*Apiurl, error) {
	ret := Apiurl{}
	err := tx.Where("code = ?", code).First(&ret).Error
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (dao *apiurlDaoImpl) Create(tx *gorm.DB, apiurl *Apiurl) (int64, error) {
	err := tx.Table(TableNameApiurl).Create(apiurl).Error
	if err != nil {
		return 0, err
	}
	return apiurl.ID, nil
}

func (dao *apiurlDaoImpl) DeleteByID(tx *gorm.DB, id int64) error {
	return tx.Where("id = ?", id).Delete(&Apiurl{}).Error
}

func (dao *apiurlDaoImpl) Update(tx *gorm.DB, apiurl *Apiurl) error {
	return tx.Table(TableNameApiurl).Where("id = ?", apiurl.ID).Updates(apiurl).Error
}

func (dao *apiurlDaoImpl) UpdateFields(tx *gorm.DB, id int64, data map[string]interface{}) error {
	return tx.Table(TableNameApiurl).Where("id = ?", id).Updates(data).Error
}
