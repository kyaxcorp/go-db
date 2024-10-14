package record

type primaryKey struct {
	// this is the initial field value
	initialFieldValue interface{}
	// this is the field name from the structure
	fieldName string
	// this is the default value from gorm!
	dbDefaultValue    string
	hasDBDefaultValue bool
	fieldType         string
}
