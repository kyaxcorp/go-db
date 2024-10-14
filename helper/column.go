package helper

import (
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"sync"
)

func GetModelDBColumns(db *gorm.DB, model interface{}) []string {
	if model == nil || db == nil {
		return nil
	}

	var columns []string
	// TODO: should DEBUG FUNCTION be removed from here?
	//result, _ := db.Debug().Migrator().ColumnTypes(model)
	result, _ := db.Migrator().ColumnTypes(model)
	for _, v := range result {
		columns = append(columns, v.Name())
	}
	return columns
}

func GetModelMapWithDBColumns(model interface{}, reverse ...bool) (map[string]string, error) {
	s, _err := schema.Parse(model, &sync.Map{}, schema.NamingStrategy{})
	if _err != nil {
		//panic("failed to create schema")
		return nil, _err
	}

	m := make(map[string]string)
	for _, field := range s.Fields {
		dbFieldName := field.DBName
		modelFieldName := field.Name
		if len(reverse) > 0 && !reverse[0] {
			m[dbFieldName] = modelFieldName
		} else {
			m[modelFieldName] = dbFieldName
		}
	}
	return m, nil
}
