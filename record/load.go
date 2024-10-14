package record

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

// ReloadData -> it will forcefully reload the data
func (r *Record) ReloadData() bool {
	r.isDbDataLoaded.False()
	return r.LoadData()
}

func (r *Record) PreparePrimaryKeys() {
	if r.primaryKeysQuery != "" {
		// it means that it's ready!
		return
	}
	var queryItems []string
	var queryItemsValues []interface{}
	//modelHelper := _struct.New(r.ModelStruct)

	// TODO: maybe all this mechanism is unnecessary because the gorm already takes by primary keys for First function

	for _, pKey := range r.primaryKeys {
		// Construct the query
		q := r.GetDBFieldName(pKey.fieldName) + " = ?"
		queryItems = append(queryItems, q)
		// Get the value for this primary key
		//queryItemsValues = append(queryItemsValues, modelHelper.GetFieldValue(pKey.fieldName))
		queryItemsValues = append(queryItemsValues, pKey.initialFieldValue)
	}
	// Create the full query
	query := strings.Join(queryItems, " AND ")

	// Save these for later usage
	r.primaryKeysQuery = query
	r.primaryKeysItemsValues = queryItemsValues
}

// LoadData -> from the db, it loads the data if exists, if it doesn't exist then it return true!
func (r *Record) LoadData() (status bool) {
	if r.isDbDataLoaded.IfFalseSetTrue() {
		return r.lastLoadDataStatus.Get()
	}

	defer func() {
		// Set the last status
		r.lastLoadDataStatus.Set(status)
	}()

	// Load the current model if exists from the db
	_db := r.getDB()
	// Get the primary key fields (names)
	// Copy the values from there to this model

	// Now based on the primary keys and based on the r.NonPtrObj we will make the query to the DB

	// The problem is that when getting the value we should transform it... there can be different types
	// so we should be careful with that! or gorm will automatically take care of it, we just need to set the raw value of it

	/*if r.IsCreateMode() {
		// the initial data doesn't have all necessary primary keys for loading the data!
		return false
	}*/

	r.PreparePrimaryKeys()

	// Set the query and its params and query
	//dbResult := _db.Where(query, queryItemsValues...).Scopes(r.getDefaultScopes).First(&r.dbData)

	// TODO:
	// TODO: use r.loadDataForUpdate for LOAD FOR UPDATE in SQL

	_db = _db.
		Where(r.primaryKeysQuery, r.primaryKeysItemsValues...)

	if r.loadDataForUpdate {
		_db = _db.Clauses(clause.Locking{Strength: "UPDATE"})
	}

	dbResult := _db.Scopes(r.getDefaultScopes).
		Take(r.dbData)
	//dbResult.RowsAffected // returns count of records found
	//dbResult.Error        // returns error or nil

	// check error ErrRecordNotFound

	if dbResult.Error != nil {
		// We have a DB error, let's check it out!
		if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
			//r.setDBError(dbResult.Error)
			r.isRecordExists.False()
			// record doesn't exist, but the loading procedure is ok!
			return true
		} else {
			// it's other type of error!
			r.setDBError(dbResult.Error)
			return false
		}
	}

	r.callOnAfterDbDataLoaded()
	r.isRecordExists.True()
	return true
}

func (r *Record) IsLoadDataOk() bool {
	return r.lastLoadDataStatus.Get()
}

func (r *Record) IsRecordNotFound() bool {
	return !r.isRecordExists.Get()
}

func (r *Record) IsRecordFound() bool {
	return r.isRecordExists.Get()
}

func (r *Record) IsDBError() bool {
	if r.GetLastDBError() != nil {
		return true
	}
	return false
}

func (r *Record) GetLoadedData() interface{} {
	return r.dbData
}
