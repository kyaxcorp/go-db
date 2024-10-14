package data_types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type IpAddresses []string

//---------------------------SQL HELPERS-----------------------\\

// https://gorm.io/docs/data_types.html
// https://pkg.go.dev/database/sql#Scanner

// Scan -> this is the func which is called when it's necessary to read from the Database to
// the defined variable with this type
func (d *IpAddresses) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	// Define a temporary variable
	var data IpAddresses
	// JSON Decode into the temporary variable
	err := json.Unmarshal(bytes, &data)
	// Check if it's ok...
	if err == nil {
		// if ok, set back the obtained value
		*d = data
	} else {
		// if error
		// empty the variable
		*d = nil
	}
	return err
}

// https://pkg.go.dev/database/sql/driver#Valuer

// Value -> this is the func which is called when it's necessary to send to Database,
// the value should be automatically converted to JSON
func (d IpAddresses) Value() (driver.Value, error) {
	// Check if there are any addresses
	if len(d) == 0 {
		// Return null/nil
		return nil, nil
	}

	// Convert the data into Json Bytes
	return json.Marshal(d)
}

//--------------------GORM HELPERS-----------------\\

// GormDataType gorm common data type
// this is the data type which is defined in the `gorm:type:DATA_TYPE` near the var!
func (IpAddresses) GormDataType() string {
	return "ip_addresses"
}

// GormDBDataType gorm db data type
// This is the Database data type which is sent to the DB
func (IpAddresses) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "JSON"
	case "postgres":
		return "JSONB"
	}
	return ""
}

//-------------- OTHER TYPE FUNCTIONS------------\\

// MarshalJSON -> convert existing data to JSON
func (d IpAddresses) MarshalJSON() ([]byte, error) {
	return json.Marshal(d)
}

// UnmarshalJSON -> set existing JSON to current data
func (d *IpAddresses) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, d)
}
