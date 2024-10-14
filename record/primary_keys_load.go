package record

import (
	"github.com/kyaxcorp/go-db/helper"
	"github.com/kyaxcorp/go-helper/_struct"
)

func (r *Record) reloadPrimaryKeys() {
	//=============== PRIMARY KEYS =================\\
	// Get the model primary keys
	primaryKeys := helper.GetModelPrimaryKeys(r.modelStruct)
	if len(primaryKeys) == 0 {
		panic("record model doesn't have any primary keys")
	}

	// create the vars
	mHelper := _struct.New(r.modelStruct)

	//mDataHelper := _struct.New(r.NonPtrObj)

	pKeys := make(map[string]primaryKey)
	var pKeysOrdered []primaryKey

	for _, pkey := range primaryKeys {
		pk := primaryKey{
			fieldName: pkey,
		}

		// set the initial field val by taking the value from the r.NonPtrObj
		//fieldValue := mDataHelper.GetFieldValue(pkey)
		fieldValue, _err := r.GetInputFieldValue(pkey)
		if _err != nil {
			// TODO: should we do anything if the field doesn't exist from the input!?
			// 		should we simply skip this field?! continue?
			// 		well, if any of the primary keys is missing it doesn't mean anything... because in this case
			// 		the user wants to create a record... and not all of the keys were supplied...
			//		so, we receive an error, and in this case we should take the value from the ModelStruct, why?
			//		because it contains the default value/state of that field!, and after that we need to use it
			//		in the find query...

			// let's get the default value
			fieldValue = _struct.GetFieldValue(r.modelStruct, pkey)
		}
		pk.initialFieldValue = fieldValue

		dbDefaultFieldVal := mHelper.GetFieldTagKeyValue(pkey, "gorm", "default")
		if dbDefaultFieldVal != "" {
			pk.dbDefaultValue = dbDefaultFieldVal
			pk.hasDBDefaultValue = true
		}

		fieldType := mHelper.GetFieldTypeName(pkey)

		pk.fieldType = fieldType

		// add to our maps...
		pKeys[pkey] = pk
		pKeysOrdered = append(pKeysOrdered, pk)
	}

	// Set to record
	r.primaryKeys = pKeys
	r.primaryKeysOrdered = pKeysOrdered
	//=============== PRIMARY KEYS =================\\
}
