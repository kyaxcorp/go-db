package filter

import (
	"github.com/google/uuid"
	"reflect"
	"time"
)

type Export struct {
	// This is a prefix for exporting/generating the file!
	ExportName string
	TableName  string
	Model      interface{}
	Columns    []ExportColumn
	Filter     *Input
	Preloads   []string
	// It will delete itself after a period of time...
	// if 0, then it will not be deleted!
	SelfDeleteAfterSeconds int64
	// items -> here we will receive the queried items

	//items []map[string]interface{}
	//items []interface{}
	items interface{}

	// if items are set
	itemsSet bool

	nrOfRows    int64
	nrOfColumns int64

	// Excel part
	excelFileID        uuid.UUID
	excelFileName      string
	excelFullFileName  string
	excelFullFilePath  string
	excelFileExtension string
	excelCreatedAt     time.Time
	excelFileSizeBytes int64
	excelError         error

	// TODO....
	pdfFileName string
}

type ExportColumn struct {
	// you can choose of these options: (it's the search criteria)
	FieldName string // it can contain .
	//DBFieldName string

	// This is the name of the column in the header
	HeaderName string
	Handler    ExportHandler
	ColWidth   float64
}

type ExportHandler func(row *ExportRow) ExportValue

type ExportValue interface{}
type ExportRow struct {
	ReflectVal reflect.Value
	RowMap     map[string]interface{}
	Row        interface{}
	FieldValue interface{}
	// Add additional functions to this structure... like GetField
}

func (e *ExportRow) GetStructValue(fieldName string) interface{} {
	if fValue, ok := e.RowMap[fieldName]; ok {
		return fValue
	}
	// TODO: should we return nil?
	return nil
}
