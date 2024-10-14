package helper

import "github.com/kyaxcorp/go-helper/_struct"

func GetModelPrimaryKeys(model interface{}) []string {
	return _struct.New(model).GetFieldNamesByTagKeyExistence("gorm", "primaryKey")
}
