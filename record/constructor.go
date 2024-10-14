package record

import (
	"reflect"
	"strings"

	"github.com/jinzhu/copier"
	"github.com/kyaxcorp/go-db/helper"
	"github.com/kyaxcorp/go-helper/Map"
	"github.com/kyaxcorp/go-helper/_context"
	"github.com/kyaxcorp/go-helper/_interface"
	"github.com/kyaxcorp/go-helper/_struct"
	"github.com/kyaxcorp/go-helper/json"
	"github.com/kyaxcorp/go-helper/ptr"
	"github.com/kyaxcorp/go-helper/sync/_bool"
)

func New(r *Record) *Record {
	if r == nil {
		r = &Record{}
	}

	//if r.NonPtrObj == nil {
	//	panic("no record model has been set")
	//}

	if r.Data == nil {
		panic("no record data has been set")
	}

	// Let's load the record from the DB
	// We will need it for comparing the values...

	//=============== CHECK IF THE MODEL IS A POINTER =================\\
	//if !_struct.IsPointer(r.NonPtrObj) {
	//	panic("the set record model is not a pointer of the inserted structure")
	//}
	//=============== CHECK IF THE MODEL IS A POINTER =================\\

	//

	// TODO: check if it's a structure or a map!

	//=============== COPY POINTER'S STRUCTURE =================\\
	//model := r.NonPtrObj
	//r.ModelStruct = &r.NonPtrObj
	//r.ModelStruct = _struct.GetPointerStructValue(r.NonPtrObj)
	//=============== COPY POINTER'S STRUCTURE =================\\

	// Check if ModelStruct has been set!

	//=============== CHECK MODEL STRUCT =================\\
	if !_struct.IsPointer(r.ModelStruct) {
		panic("model struct is not a pointer of model")
	}
	if r.ModelStruct == nil {
		panic("no model struct has been set")
	}
	// Get the structure!
	//r.modelStruct = _struct.GetPointerStructValue(r.ModelStruct)
	r.modelStruct = _interface.CloneInterfaceItem(r.ModelStruct)
	//=============== CHECK MODEL STRUCT =================\\

	//

	// Create the model structure
	r.dataCopied = _interface.CloneInterfaceItem(r.modelStruct)
	r.inputFieldNames = make(map[string]string)
	r.inputData = make(map[string]interface{})

	//=============== CHECK IF DATA IS A POINTER =================\\
	if Map.IsMap(r.Data) {
		// Check if it's map[string]interface{}

		mapType := reflect.TypeOf(r.Data).String()
		if mapType != "map[string]interface{}" && mapType != "map[string]interface {}" {
			panic("map is not of type map[string]interface{}, it's -> " + mapType)
		}

		r.dataMap = r.Data.(map[string]interface{})

		if r.dataMap == nil {
			panic("data map is nil")
		}

		// Check if there are any data...
		if len(r.dataMap) == 0 {
			panic("data map has no items")
		}

		// Let's check the ID field...

		if val, ok := r.dataMap["ID"]; ok {
			//idType := reflect.TypeOf(val)
			idValue := reflect.ValueOf(val)
			if idValue.IsZero() {
				// Then we should remove it!
				delete(r.dataMap, "ID")
			}
		}

		// let's copy each value from map to Struct!

		var _err error
		// Now we should convert the incoming values to the model's format!
		// But we should select only the values that have been set!
		dataMap2, _err := Map.ConvertMapValuesBasedOnModel(r.dataMap, r.modelStruct, nil)
		if _err != nil {
			panic("failed to convert map values based on model -> r.dataMap, r.modelStruct -> " + _err.Error())
		}
		r.dataMap = Map.RecursiveCopyToExistent(dataMap2, r.dataMap)

		// TODO: copier right now can't copy map[string]interface{} to Structure...
		// Maybe someday later it may have this functionality!
		// That's why we will be using json encode/decode

		r.dataMapJson, _err = json.Encode(r.dataMap)
		if _err != nil {
			panic("failed to convert r.dataMap to r.dataMapJson -> " + _err.Error())
		}

		// Copy the incoming data (json formatted) to r.dataCopied which is a clone of the Model!
		_err = json.Decode(r.dataMapJson, r.dataCopied)
		if _err != nil {
			panic("failed to convert r.dataMapJson r.dataCopied -> " + _err.Error())
		}

		// Now let's Get a map of the copied model data structure
		tmpDataMap := _struct.New(r.dataCopied).Map()

		// Let's now copy the incoming data (that is formatted) to r.inputData
		for fieldName, _ := range r.dataMap {
			r.inputFieldNames[fieldName] = ""
			r.inputData[fieldName] = tmpDataMap[fieldName]
		}

		//_err := copier.Copy(r.dataCopied, r.dataMap)
		//if _err != nil {
		//	panic("failed to copy data from dataStr to dataCopied -> " + _err.Error())
		//}

		// Set the type of the received data
		r.inputDataType = inputDataMapInterface
	} else {
		// check if it's a pointer
		if !ptr.Is(r.Data) {
			// Check if it's a struct
			panic("data is not a pointer")
		}

		// Check if it's a struct!
		if !_struct.IsStruct(r.Data) {
			panic("data is not a structure")
		}

		// Get the data from the pointer
		//r.dataStr = _struct.GetPointerStructValue(r.Data)
		r.dataStr = _interface.CloneInterfaceItem(r.Data)

		_err := copier.Copy(r.dataCopied, r.dataStr)
		if _err != nil {
			panic("failed to copy data from dataStr to dataCopied -> " + _err.Error())
		}

		tmpDataMap := _struct.New(r.dataCopied).Map()

		// Let's set the field names
		tmpStr := _struct.New(r.dataStr).Map()
		for fieldName, _ := range tmpStr {
			r.inputFieldNames[fieldName] = ""

			r.inputData[fieldName] = tmpDataMap[fieldName]
		}

		// Set the received type
		r.inputDataType = inputDataStruct
	}

	if val, ok := r.inputData["ID"]; ok {
		//idType := reflect.TypeOf(val)
		idValue := reflect.ValueOf(val)
		if idValue.IsZero() {
			// Then we should remove it!
			delete(r.inputData, "ID")
		}
	}

	//var _err error
	//r.dataMap, _err = Map.ConvertMapValuesBasedOnModel(r.dataMap, r.modelStruct, nil)
	//if _err != nil {
	//	panic("failed to convert map values based on model -> r.dataMap, r.modelStruct -> " + _err.Error())
	//}

	// Set the model struct helper
	r.dataStrHelper = _struct.New(r.dataCopied)

	//=============== CHECK IF DATA IS A POINTER =================\\

	//

	//=============== CHECK IF DB CLIENT HAS BEEN SET =================\\
	// if r.DB == nil {
	// 	// Get the default one!
	// 	_db, _err := db.GetDefaultClient()
	// 	if _err != nil {
	// 		panic(_err.Error())
	// 	}
	// 	r.DB = _db
	// }
	//=============== CHECK IF DB CLIENT HAS BEEN SET =================\\

	//

	//=============== CHECK CONTEXT HAS BEEN SET =================\\
	if r.Ctx == nil {
		r.Ctx = _context.GetDefaultContext()
	}
	//=============== CHECK CONTEXT HAS BEEN SET =================\\

	//

	//=============== CREATE NECESSARY VARS =================\\

	r.dbSet = _bool.New()
	r.cachedIDSet = _bool.New()
	r.saveModeSet = _bool.New()
	r.isDbDataLoaded = _bool.New()
	r.lastLoadDataStatus = _bool.New()
	r.isRecordExists = _bool.New()
	r.disableDefaultScope = _bool.NewVal(false)

	// Copy the inserted model to the dbData (we need to have the same structure!)
	//r.dbData = r.NonPtrObj

	// Copy the ModelStruct to the dbData for Gorm to know where to save the loaded data!
	r.dbData = _interface.CloneInterfaceItem(r.modelStruct)
	//=============== CREATE NECESSARY VARS =================\\

	//

	//=============== PRIMARY KEYS =================\\
	r.reloadPrimaryKeys()
	//=============== PRIMARY KEYS =================\\

	//

	//=============== OVERRIDE SAVE MODE =================\\

	if r.SaveMode != "" {
		r.SaveMode = strings.ToLower(r.SaveMode)
	}
	//=============== OVERRIDE SAVE MODE =================\\

	//

	//=============== GET DB COLUMN NAMES=================\\
	m, _err := helper.GetModelMapWithDBColumns(r.modelStruct)
	if _err != nil {
		panic(_err.Error())
	}
	r.modelFieldNamings = m
	//=============== GET DB COLUMN NAMES=================\\

	// Check the mode?!
	// If the ID is empty... then it's in create mode, if it's not empty then it's in save mode

	//=============== LOAD THE DATA FOR EXISTING RECORD =================\\
	r.PreparePrimaryKeys()
	if r.IsSaveMode() {
		if r.AutoLoad && !r.LoadData() {
			panic(r.GetLastDBError().Error())
		}
	}
	//=============== LOAD THE DATA FOR EXISTING RECORD =================\\

	return r
}

func GetInterfacePointer(obj interface{}) interface{} {
	p := reflect.New(reflect.TypeOf(obj))
	p.Elem().Set(reflect.ValueOf(obj))
	return p.Interface()
}
