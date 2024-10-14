package helper

import (
	"gorm.io/gorm"
	"time"
)

func DeletedAt(t *time.Time) gorm.DeletedAt {
	deletedAt := gorm.DeletedAt{}
	if t != nil {
		deletedAt.Time = *t
		deletedAt.Valid = true
	}
	return deletedAt
}

func DeletedAtP(t *time.Time) *gorm.DeletedAt {
	d := DeletedAt(t)
	return &d
}

func DeletedAtNow() gorm.DeletedAt {
	deletedAt := gorm.DeletedAt{}
	deletedAt.Time = time.Now()
	deletedAt.Valid = true
	return deletedAt
}

func DeletedAtNowP() *gorm.DeletedAt {
	d := DeletedAtNow()
	return &d
}
