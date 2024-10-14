package record

import (
	"github.com/google/uuid"
	"reflect"
)

const (
	ModeSave   = "save"
	ModeCreate = "create"
)

func (r *Record) getSaveMode() string {

	if r.saveModeSet.IfFalseSetTrue() {
		return r.SaveMode
	}

	if r.SaveMode != "" {
		switch r.SaveMode {
		case "save", "update":
			return ModeSave
		case "create", "insert":
			return ModeCreate
		default:
			// Should be detected...
		}
	}
	// then use the autodetection

	// We should check if the field ID is defined, but there may be other primary keys...
	// anyway, if some of the primary keys are empty... it means that it's create, but that's not always true!
	// because some primary keys may have gen_default_uuidv4 or something like that..., so the value it's been generated
	// by the database itself!
	// so if we don't know that operation should be made... it's better to have an override value which specifies
	// what is it! if the auto detection doesn't work as needed, then use the override...

	// Check if the field ID exists...
	// Check if it has value...
	// Check the other primary keys...

	/*
		We should know more information about the primary keys like:
		- if it has a default value set in the DB, if yes, then it's not mandatory to be set here in the NonPtrObj!


		But how do we detect that is an Insert?! Well, very easily...
		Gorm cannot know what do you want from it... so we should make our own algo!

		1. If all values (primary keys) are been set in the model, and some of the primary keys has default values in the DB,
		it means that it's a save... usually it should be like that! Because in that case we will let the DB to set the default value...
		But sometimes we would like to override these values even if there are some default, in this case we need to have
		an option that tells that we want to override... or simply indicate the MODE! so we are talking here about auto detection!
		2. If the primary keys don't have a default value set for the DB! then what?! How do we know what to choose?!
		In this case we should respond that this is a SAVE MODE, but if the user needs to create, then it should set the override Mode to Create!


	*/

	// TODO: the mode should be cached...

	isEmpty := false
	for _, pKey := range r.primaryKeysOrdered {

		// by default it should be for SAVE

		// if not set value and there is a default db value, then it's create
		// if not set value and there is no default db value, then it's error!
		// if set value and there is default db value, then it's save
		// if set value and there is no default db value, then it's for save...

		//log.Println(pKey.fieldName, pKey.initialFieldValue)

		switch pKey.fieldType {
		case "UUID":
			if pKey.initialFieldValue.(uuid.UUID) == uuid.Nil {
				isEmpty = true
				//log.Println(pKey.fieldName, "is nil")
			}
		/*case "Time":

		case "string":
			if pKey.initialFieldValue.(string) == "" {
				isEmpty = true
			}
		case "bool":
			// cannot detect...*/
		default:
			if reflect.ValueOf(pKey.initialFieldValue).IsZero() {
				isEmpty = true
				//log.Println(pKey.fieldName, "is zero")
			}
		}

		// if it's empty... then check if there is a DB default value set...

		if isEmpty {
			if pKey.hasDBDefaultValue {
				// It's CREATE MODE!
				break
			} else {
				// reset..., check the next one
				isEmpty = false
			}
		}
	}

	// if it's empty... it means that we should set to create!
	mode := ModeSave
	if isEmpty {
		mode = ModeCreate
	}
	r.SaveMode = mode

	return mode
}

func (r *Record) IsCreateMode() bool {
	if r.getSaveMode() == ModeCreate {
		return true
	}
	return false
}

func (r *Record) IsSaveMode() bool {
	if r.getSaveMode() == ModeSave {
		return true
	}
	return false
}
