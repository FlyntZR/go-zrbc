package structUtil

import (
	"reflect"

	"github.com/pkg/errors"
)

//全局变量

// 结构体tag，标记为时间格式 在tag中增加 time_format:"true"
const IS_TIME_FORMAT_TAG = "time_format"

// 有无用作标题的tag
const HAS_TITLE_TAG = "title"

// 标记为不导出到excel
const SKIP_TAG = "skip"
const DESCRIPTION_TAG = "description"
const GORM_TAG = "gorm"
const JSON_TAG = "json"
const DEMO_TAG = "demo"
const TIME_FORMAT_DAY = "2006/01/02"
const TIME_FORMAT_MINUTE = "2006/01/02 15:04:05"

const MSG_TYPE_MINIPROGRAM = "miniprogram_notice"

// 结构体描述中的gorm-json-description
type StructField struct {
	Column       string              `description:"gorm信息"`
	Json         string              `description:"字段的json描述"`
	Description  string              `description:"字段的名称，注释/描述"`
	IsTimeFormat string              `description:"是否为时间格式"`
	Demo         string              `description:"样例值"`
	TypeField    reflect.StructField `description:"类型"`
	ValueField   reflect.Value       `description:"值"`
}

// 利用反射获取结构体标签信息
// 传入结构体，返回结构体字段序号和对应Tag标签值的键值对（标签值为json描述中的description）
func ParsingStruct(toParse interface{}) (map[int]StructField, error) {
	typeOfStruct := reflect.TypeOf(toParse)
	valueOfStruct := reflect.ValueOf(toParse)

	kind := valueOfStruct.Kind()
	if kind != reflect.Struct {
		return nil, errors.New("ParsingStruct Need Struct")
	}
	idTagMap := make(map[int]StructField, 0)
	numOfField := valueOfStruct.NumField()
	for id := 0; id < numOfField; id++ {
		var des StructField
		des.Column = typeOfStruct.Field(id).Tag.Get(GORM_TAG)
		des.Json = typeOfStruct.Field(id).Tag.Get(JSON_TAG)
		des.Description = typeOfStruct.Field(id).Tag.Get(DESCRIPTION_TAG)
		des.IsTimeFormat = typeOfStruct.Field(id).Tag.Get(IS_TIME_FORMAT_TAG)
		des.Demo = typeOfStruct.Field(id).Tag.Get(DEMO_TAG)
		des.TypeField = typeOfStruct.Field(id)
		des.ValueField = valueOfStruct.Field(id)
		idTagMap[id] = des
	}
	return idTagMap, nil
}
