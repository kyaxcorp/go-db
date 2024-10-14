package scope

import "gorm.io/gorm"

func OrderByCreatedAtAsc(db *gorm.DB) *gorm.DB {
	return db.Order("created_at asc")
}

func OrderByCreatedAtDesc(db *gorm.DB) *gorm.DB {
	return db.Order("created_at desc")
}

func OrderByUpdatedAtAsc(db *gorm.DB) *gorm.DB {
	return db.Order("updated_at asc")
}

func OrderByUpdatedAtDesc(db *gorm.DB) *gorm.DB {
	return db.Order("updated_at desc")
}

func OrderByDeletedAtAsc(db *gorm.DB) *gorm.DB {
	return db.Order("deleted_at asc")
}

func OrderByDeletedAtDesc(db *gorm.DB) *gorm.DB {
	return db.Order("deleted_at desc")
}
