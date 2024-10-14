package record

import (
	"time"

	"github.com/kyaxcorp/go-helper/_struct"
	"github.com/kyaxcorp/go-helper/errors2"
	"gorm.io/gorm"
)

func (r *Record) Delete() bool {
	_db := r.getDB()

	if !r.prepareSaveData() {
		return false
	}

	var _err error
	uID := r.GetUserID()
	uIDisNil := r.isUserIDNil()

	_err = r.callOnBeforeDelete()
	if _err != nil {
		r.setError(_err)
		r.callOnError()
		r.callOnSaveError()
		return false
	}

	var result *gorm.DB
	if r.IsCreateMode() {
		_err := errors2.New(0, "record doesn't exist or id not set")
		r.setDBError(_err)
		return false
	} else {
		// We should update it!
		// TODO: should we get the record before making changes?!...

		isDeletedExists := _struct.FieldExists(r.modelStruct, "IsDeleted")
		deletedAtExists := _struct.FieldExists(r.modelStruct, "DeletedAt")

		if isDeletedExists || deletedAtExists {
			//r.saveData["IsDeleted"] = true
			if isDeletedExists {
				r.SetSaveFieldValue("IsDeleted", true)
			}

			if deletedAtExists {
				//r.saveData["DeletedAt"] = time.Now()
				r.SetSaveFieldValue("DeletedAt", time.Now())
			}
			if !uIDisNil && _struct.FieldExists(r.modelStruct, "DeletedByID") {
				//r.saveData["DeletedBy"] = uID
				r.SetSaveFieldValue("DeletedByID", uID)
			}

			_err = r.callOnBeforeUpdate()
			if _err != nil {
				r.setError(_err)
				r.callOnError()
				r.callOnSaveError()
				return false
			}
			//saveDataModel := r.generateSaveDataModel()
			result = _db.Omit(r.GetOmitFields()...).Save(r.saveData)
			r.loadDataForUpdate = false
			r.dbData = r.saveData

			//r.ReloadData()
			_err = r.callOnAfterUpdate()
			if _err != nil {
				r.setError(_err)
				r.callOnError()
				r.callOnSaveError()
				return false
			}

			if result.Error != nil {
				r.setDBError(result.Error)
				return false
			}
			_err = r.callOnAfterDelete()
			if _err != nil {
				r.setError(_err)
				r.callOnError()
				r.callOnSaveError()
				return false
			}
			return true
		} else {
			return r.ForceDelete()
		}
	}
}

// ForceDelete -> it will completely delete the record from the DB
func (r *Record) ForceDelete() bool {
	_db := r.getDB()
	// TODO: should we get the record before making changes?!...
	// TODO: should we check if exists in the DB?!...or gorm does it...

	var _err error
	_err = r.callOnBeforeForceDelete()
	if _err != nil {
		r.setError(_err)
		r.callOnError()
		r.callOnSaveError()
		return false
	}

	var result *gorm.DB
	if r.IsCreateMode() {
		_err := errors2.New(0, "record doesn't exist or id not set")
		r.setDBError(_err)
		r.callOnError()
		r.callOnDeleteError()
		return false
	} else {
		// We should delete it!
		// Does gorm know the primary keys?! when there are multiple
		result = _db.Delete(r.saveData)
		if result.Error != nil {
			r.setDBError(result.Error)
			r.callOnError()
			r.callOnDBError()
			r.callOnDeleteError()
			return false
		}
		_err = r.callOnAfterForceDelete()
		if _err != nil {
			r.setError(_err)
			r.callOnError()
			r.callOnSaveError()
			return false
		}
		return true
	}
	// TODO: add hooks
}
