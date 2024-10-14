package scope

import (
	"gorm.io/gorm"
)

func IsDeleted(db *gorm.DB) *gorm.DB {
	return db.Where("is_deleted = ?", true)
}

func IsDeletedAt(db *gorm.DB) *gorm.DB {
	return db.Where("deleted_at IS NOT NULL")
}

func IsDeletedT(tableName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(tableName+".is_deleted = ?", true)
	}
}

func IsDeletedAtT(tableName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(tableName + ".deleted_at IS NOT NULL")
	}
}

func IsNotDeleted(db *gorm.DB) *gorm.DB {
	// Migrator has column queries the DB information schema for existence,
	// but this is slow in performance, that's why, we will check the struct itself (the provided model)
	/*if db.Migrator().HasColumn(db.Statement.Model, "IsDeleted") {
		return db.Where("is_deleted = ?", false)
	} else if db.Migrator().HasColumn(db.Statement.Model, "DeletedAt") {
		return db.Where("deleted_at IS NULL")
	}*/

	// din pacate db.Statement.Model e nil... ceea ce inseamna ca nu putem sti despre ce model/tabel e vorba
	/*model := _struct.New(db.Statement.Model)
	if model.FieldExists("IsDeleted") {
		return db.Where("is_deleted = ?", false)
	} else if model.FieldExists("DeletedAt") {
		return db.Where("deleted_at IS NULL")
	}
	return db
	*/
	return db.Where("is_deleted = ?", false)
}
func IsNotDeletedAt(db *gorm.DB) *gorm.DB {
	return db.Where("deleted_at IS NULL")
}

func IsNotDeletedT(tableName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(tableName+".is_deleted = ?", false)
	}
}

func IsNotDeletedAtT(tableName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(tableName + ".deleted_at IS NULL")
	}
}

func IsActive(db *gorm.DB) *gorm.DB {
	return db.Where("is_active = ?", true)
}

func IsActiveT(tableName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(tableName+".is_active = ?", true)
	}
}

func IsNotActive(db *gorm.DB) *gorm.DB {
	return db.Where("is_active = ?", false)
}

func IsNotActiveT(tableName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(tableName+".is_active = ?", false)
	}
}

func IsExpired(db *gorm.DB) *gorm.DB {
	return db.Where("is_expired = ?", true)
}

func IsExpiredT(tableName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(tableName+".is_expired = ?", true)
	}
}

func IsNotExpired(db *gorm.DB) *gorm.DB {
	return db.Where("is_expired = ?", false)
}

func IsNotExpiredT(tableName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(tableName+".is_expired = ?", false)
	}
}
