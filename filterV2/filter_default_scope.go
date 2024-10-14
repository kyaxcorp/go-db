package filterV2

import (
	"github.com/kyaxcorp/go-db/scope"
	"gorm.io/gorm"
)

func (f *Input) EnableDefaultScope() *Input {
	f.enableDefaultScope = true
	return f
}

func (f *Input) DisableDefaultScope() *Input {
	f.enableDefaultScope = false
	return f
}

func (f *Input) applyDefaultScope() *Input {
	f.db = f.db.Scopes(f.getDefaultScope)
	return f
}

func (f *Input) getDefaultScope(db *gorm.DB) *gorm.DB {
	// Check if field exists! (check all models?!)
	// We should check the model?!
	// TODO: should we cache here anything?!

	var _err error
	if _, _err = f.getDBFieldName("IsDeleted"); _err == nil {
		//db.Scopes(scope.IsNotDeleted)
		// TODO: how we should do here?!
		db.Scopes(scope.IsNotDeletedT(f.primaryModel.dbTableName))
	} else if _, _err = f.getDBFieldName("DeletedAt"); _err == nil {
		db.Scopes(scope.IsNotDeletedAtT(f.primaryModel.dbTableName))
	}

	return db
}
