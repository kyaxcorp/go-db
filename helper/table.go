package helper

import "gorm.io/gorm"

func GetModelDBTableName(db *gorm.DB, model interface{}) (string, error) {
	stmt := &gorm.Statement{DB: db}
	_err := stmt.Parse(model)
	if _err != nil {
		//panic("failed to parse model structure")
		return "", _err
	}
	return stmt.Schema.Table, nil
}
