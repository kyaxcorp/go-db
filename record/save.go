package record

import (
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/kyaxcorp/go-helper/_interface"
	"github.com/kyaxcorp/go-helper/_struct"
	"github.com/kyaxcorp/go-helper/errors2/define"
	"github.com/kyaxcorp/go-helper/json"
	"gorm.io/gorm"
)

/*
Create some kind of structure which will be embedded in another...?!
It will take the functions, but it will not take the values, because it will not know the parent structure!
Only if we create a function which sets the reference to the parent... in this case it may work...

*/

// TODO: we can set global hooks for this structure...

func (r *Record) getInputOmitFields() []string {
	_strMap := _struct.New(r.modelStruct).Map()
	var omitFields []string

	// TODO: should we also omit for nested structures that have not been set?!
	// TODO: or we should simply load their data and add to it...

	// gorm allows to omit sub fields! they should be declared with . (dots)

	for fieldName, _ := range _strMap {
		//if _, fieldFound := r.saveData[fieldName]; !fieldFound {
		if _, fieldFound := r.dataMap[fieldName]; !fieldFound {
			// if something is not present, then omit it
			omitFields = append(omitFields, fieldName)
		}
	}
	return omitFields
}

// func (r *Record) GetSaveData() map[string]interface{} {
func (r *Record) GetSaveData() interface{} {
	return r.saveData
}

// TODO: get save data field and set save data field!

// func (r *Record) SetSaveData(saveData map[string]interface{}) *Record {
func (r *Record) SetSaveData(saveData interface{}) *Record {
	r.saveData = saveData
	return r
}

func (r *Record) SetOmitField(fieldName string) {
	//  TODO: Check if field already present!!!
	//r.omitFields = append(r.omitFields, fieldName)
	r.omitFields[fieldName] = nil
}

func (r *Record) RemoveOmitField(fieldName string) {
	//  TODO: Check if field already present!!!
	//r.omitFields = append(r.omitFields, fieldName)
	if _, ok := r.omitFields[fieldName]; ok {
		delete(r.omitFields, fieldName)
	}
}

func (r *Record) GetOmitFields() []string {
	var f []string
	for fieldName, _ := range r.omitFields {
		f = append(f, fieldName)
	}
	return f
}

func (r *Record) SetSaveFieldValue(fieldName string, value interface{}) {
	// Check if there is a set field already
	// SetOmitField
	r.RemoveOmitField(fieldName)
	_struct.New(r.saveData).SetInterface(fieldName, value)
	//reflect.ValueOf(r.saveData).Elem().FieldByName(fieldName).Set(reflect.ValueOf(value))
}

func (r *Record) generateSaveDataModel() interface{} {
	_model := _interface.CloneInterfaceItem(r.modelStruct)
	_json, _err := json.Encode(r.saveData)
	if _err != nil {
		panic("failed to encode r.saveData to Json -> " + _err.Error())
	}
	_err = json.Decode(_json, _model)
	if _err != nil {
		panic("failed to Decode _json to _model -> " + _err.Error())
	}
	return _model
}

func (r *Record) Save() bool {
	_db := r.getDB()

	var result *gorm.DB
	var _err error

	if !r.prepareSaveData() {
		return false
	}

	_err = r.callOnBeforeSave()
	if _err != nil {
		r.setError(_err)
		r.callOnError()
		r.callOnSaveError()
		return false
	}

	uID := r.GetUserID()
	uIDisNil := r.isUserIDNil()

	// TODO: la noi ID e uuid.UUID insa in input el vine ca string

	// Updated should be always present!
	if !uIDisNil && _struct.FieldExists(r.modelStruct, "UpdatedByID") {
		//r.saveData["UpdatedBy"] = uID
		r.SetSaveFieldValue("UpdatedByID", uID)
	}
	//
	if _struct.FieldExists(r.modelStruct, "UpdatedAt") {
		//r.saveData["UpdatedAt"] = time.Now()
		r.SetSaveFieldValue("UpdatedAt", time.Now())
	}

	if r.IsCreateMode() {
		// If it's nil, then we should create it!
		if !uIDisNil && _struct.FieldExists(r.modelStruct, "CreatedByID") {
			// check which type is user id -> uuid or other type
			//r.saveData["CreatedBy"] = uID
			r.SetSaveFieldValue("CreatedByID", uID)
		}
		//
		if _struct.FieldExists(r.modelStruct, "CreatedAt") {
			//r.saveData["CreatedAt"] = time.Now()
			r.SetSaveFieldValue("CreatedAt", time.Now())
		}

		// 1. copy the data to the real structure
		// 2. omit all fields that are not in the list (we should not include primary keys in the omit list!) because we should receive them back!
		// 3. Create
		// 4. Now let's read the data only by list and set it to dbData!
		// 5. launch reload to load the entire data!

		//saveDataModel := r.generateSaveDataModel()
		result = _db.Omit(r.GetOmitFields()...).Create(r.saveData)
		r.loadDataForUpdate = false
		r.dbData = r.saveData
		// TODO: later on we should do a reload like on save?!

		_err = r.callOnAfterInsert()
		if _err != nil {
			r.setError(_err)
			r.callOnError()
			r.callOnSaveError()
			return false
		}
	} else {
		_err = r.callOnBeforeUpdate()
		if _err != nil {
			r.setError(_err)
			r.callOnError()
			r.callOnSaveError()
			return false
		}
		//saveDataModel := r.generateSaveDataModel()
		result = _db.Omit(r.GetOmitFields()...).Save(r.saveData)
		r.dbData = r.saveData
		r.ReloadData()

		// We should update it!
		//result = _db.Save(r.saveData)
		_err = r.callOnAfterUpdate()
		if _err != nil {
			r.setError(_err)
			r.callOnError()
			r.callOnSaveError()
			return false
		}
	}
	_err = r.callOnAfterSave()
	if _err != nil {
		r.setError(_err)
		r.callOnError()
		r.callOnSaveError()
		return false
	}

	// We can return BOOL and save the error somewhere in the record as the last error!

	if result.Error != nil {
		r.setDBError(result.Error)
		r.callOnError()
		r.callOnDBError()
		r.callOnSaveError()
		return false
	}

	// Load back the data!

	return true
}

func (r *Record) isUserIDNil() bool {
	uID := r.GetUserID()
	uIDv := reflect.ValueOf(uID)
	uIDType := uIDv.Type().String()
	uIDisNil := false
	if uIDType == "uuid.UUID" || uIDType == "*uuid.UUID" {
		// TODO: we should test if it's nil by using pointer or not
		if uID == uuid.Nil {
			uIDisNil = true
		}
	} else {
		if uID == nil {
			uIDisNil = true
		}
	}
	return uIDisNil
}

func (r *Record) GetSavedData() interface{} {
	//return r.saveData
	return r.dbData
}

func (r *Record) prepareSaveData() bool {
	callSimpleError := func(_err error) {
		r.setError(_err)
		r.callOnError()
		r.callOnSaveError()
	}

	// let's recreate the map with no keys...
	//r.saveData = make(map[string]interface{})

	//var dbData interface{}

	if r.IsSaveMode() {
		// if it's save mode then we should get the loaded data from r.dbData
		// and after that put over it the inputData

		r.loadDataForUpdate = true
		r.ReloadData()
		r.saveData = _interface.CloneInterfaceItem(r.dbData)
		//dbDataMap := _struct.New(r.dbData).Map()
		//r.dataMap
		// copy first the current db data to saveData
		//Map.CopyStringInterface(dbDataMap, r.saveData)
	} else {
		r.saveData = _interface.CloneInterfaceItem(r.modelStruct)
	}

	_err := json.Decode(r.dataMapJson, r.saveData)
	if _err != nil {
		panic("failed to convert r.dataMapJson r.saveData -> " + _err.Error())
	}

	r.omitFields = make(map[string]*bool)
	omitFields := r.getInputOmitFields()
	for _, fieldName := range omitFields {
		r.SetOmitField(fieldName)
	}

	// let's copy the inputData to save data, why? because we don't want to flood the inputData with other information
	// the saveData variable can have or can be supplied with other additional information!
	//Map.CopyStringInterface(r.inputData, r.saveData)
	//Map.CopyStringInterface(r.dataMap, r.saveData)

	// step 3 - copy the data from the input
	switch r.inputDataType {
	case inputDataMapInterface:

	case inputDataStruct:

	default:
		callSimpleError(define.Err(0, "unknown input data"))
		return false
	}
	// Check some fields if the types are correct....
	return true
}
