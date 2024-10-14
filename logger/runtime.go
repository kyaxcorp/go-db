package logger

import (
	loggerConfig "github.com/kyaxcorp/go-logger/config"
	"gorm.io/gorm"
)

// Create a runtime logger

type Runtime struct {
	InstanceName string
	DBType       string
	Level        int
}

func NewQuickConsole(db *gorm.DB, r Runtime) error {
	level := 1
	if r.Level > 0 {
		level = r.Level
	}

	l := NewLogger(
		r.InstanceName,
		r.DBType,
		loggerConfig.Config{
			Level: level,
		},
	)
	db.Logger = l
	return nil
}

func QuickScope(r Runtime) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		dbc := db.Session(&gorm.Session{})
		NewQuickConsole(dbc, r)
		return dbc
	}
}
