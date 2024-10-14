package filter

import (
	//"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/kyaxcorp/go-helper/_struct"
	"github.com/kyaxcorp/go-helper/conv"
	"github.com/kyaxcorp/go-helper/export"
	"github.com/kyaxcorp/go-helper/function"
	"github.com/xuri/excelize/v2"

	//"github.com/kyaxcorp/go-helper/err/define"
	"reflect"

	"github.com/kyaxcorp/go-helper/file"
	"github.com/kyaxcorp/go-helper/filesystem"
	"gorm.io/gorm"

	"time"
)

func (e *Export) GeneratePdf() {

}

func (e *Export) GetExcelExportPath() (string, error) {
	return export.GetExcelFileExportPath()
}

func (e *Export) GetPdfExportPath() (string, error) {
	return export.GetPdfFileExportPath()
}

func (e *Export) GetHtmlExportPath() (string, error) {
	return export.GetHtmlFileExportPath()
}

func (e *Export) GetJsonExportPath() (string, error) {
	return export.GetJsonFileExportPath()
}

func (e *Export) GetWordExportPath() (string, error) {
	return export.GetWordFileExportPath()
}

func (e *Export) GenerateExcel() bool {
	// based on the input Filter, we should query and generate

	if e.nrOfRows == 0 {
		return false
	}

	// Let's initially export everything...
	f := excelize.NewFile()

	sheetName := "Sheet1"
	sheetIndex, _err := f.GetSheetIndex(sheetName)
	if _err != nil {
		e.excelError = _err
		return false
	}
	// First, let's create the struct for speeding up the process!

	type HeaderField struct {
		Header    string
		FieldName string
		//DBFieldName string
		Handler ExportHandler
		XAxis   string // this is the column
		// A,B,C etc...
		Column   string
		ColWidth float64
	}

	//firstRow := e.items[0]
	// TODO: later on if no columns defined, we can set to get from the DB or model/table

	//var mFields = make(map[string]string)
	//if e.Model != nil {
	//	mFields, _err = dbHelper.GetModelMapWithDBColumns(e.Model, true)
	//	if _err != nil {
	//		e.excelError = _err
	//		return false
	//	}
	//}

	var headerFields []HeaderField
	for columnNr, columnDetails := range e.Columns {
		clNr := columnNr + 1

		headerName := columnDetails.HeaderName
		//dbColumn := ""
		fieldName := ""

		var _err error
		if columnDetails.FieldName != "" {
			fieldName = columnDetails.FieldName

			// Check if this field exists!

			//if !_struct.FieldExists(e.Model, fieldName) {
			//	log.Println("field doesn't exist!!!", fieldName)
			//	continue
			//}

			//dbColumn, _err = e.Filter.getDBFieldName(columnDetails.FieldName)
			//if _err != nil {
			//	e.excelError = _err
			//	return false
			//}
		} else if function.IsCallable(columnDetails.Handler) {
			// if there is a callback
		} else {
			continue
		}
		//else if columnDetails.DBFieldName != "" {
		//	dbColumn = columnDetails.DBFieldName
		//	if fName, ok := mFields[dbColumn]; ok {
		//		fieldName = fName
		//	} else {
		//		e.excelError = define.Err(0, "structure field not found from database field name -> ", dbColumn)
		//		return false
		//	}
		//}

		if headerName == "" {
			if columnDetails.FieldName != "" {
				headerName = columnDetails.FieldName
			}
			//else {
			//	headerName = columnDetails.DBFieldName
			//}
		}

		xAxis, _err := excelize.ColumnNumberToName(clNr)
		if _err != nil {
			e.excelError = _err
			return false
		}

		headerFields = append(headerFields, HeaderField{
			Header:    headerName,
			FieldName: fieldName,
			//DBFieldName: dbColumn,
			Handler:  columnDetails.Handler,
			XAxis:    xAxis,
			Column:   xAxis,
			ColWidth: columnDetails.ColWidth,
		})
	}

	for _, headerField := range headerFields {
		excelRowNr := 1

		/*style, err := f.NewStyle(
			`{
					"alignment":{
						"horizontal":"center",
						"WrapText":true
					},
					"font":{
						"bold":true
					}
				}`,
		)*/
		style, err := f.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{
				Horizontal: "center",

				//ShrinkToFit: false,
				//WrapText:    false,
			},
			Font: &excelize.Font{
				Bold: true,
			},
		},
		)
		if err != nil {
			//fmt.Println(err)
			e.excelError = _err
			return false
		}

		XYAxis := headerField.XAxis + conv.IntToStr(excelRowNr)

		_err := f.SetCellStyle(sheetName, XYAxis, XYAxis, style)
		if _err != nil {
			e.excelError = _err
			return false
		}

		if headerField.ColWidth > 0 {
			_err = f.SetColWidth(sheetName, headerField.XAxis, headerField.XAxis, headerField.ColWidth)
			if _err != nil {
				e.excelError = _err
				return false
			}
		}

		_err = f.SetCellValue(sheetName, XYAxis, headerField.Header)
		if _err != nil {
			e.excelError = _err
			return false
		}
	}

	slice := reflect.ValueOf(e.items)
	sliceLen := slice.Len()

	for rowNr := 0; rowNr < sliceLen; rowNr++ {
		row := slice.Index(rowNr)
		rowInterface := row.Interface()
		rowMap := _struct.New(rowInterface).Map()
		exportRow := &ExportRow{
			ReflectVal: row,
			Row:        rowInterface,
			RowMap:     rowMap,
		}

		for _, headerField := range headerFields {
			excelRowNr := rowNr + 2
			// row contains the db field names!

			var fieldValue interface{}

			// we can have fields with ptr and without
			// if we see pointer value, and it's zero, should create it to Zero Value or should we set to nil!?
			// the thing is that there can be a Handler that process the value... and it can give an error because the
			// value is simply nil
			// so lets by default get the value or create a zero value in case of a pointer!

			if headerField.FieldName != "" {
				if strings.Contains(headerField.FieldName, ".") {
					//fields := strings.Split(headerField.FieldName, ".")
					fieldValueRef, _err := _struct.GetNestedFieldReflectValue(row, headerField.FieldName)
					if _err != nil {
						//fieldValue = nil
						// it's better for us to set as empty string than nil, because nil gives an error:
						// panic: runtime error: invalid memory address or nil pointer dereference
						fieldValue = ""
						//log.Println(headerField.FieldName, _err)

					} else {
						fieldValue = fieldValueRef.Interface()
						//log.Println(headerField.FieldName, fieldValue)
					}
				} else {
					fieldValue = row.FieldByName(headerField.FieldName).Interface()
				}
			}

			// ============== transformations from pointer ==============\\

			if fieldValue == nil {

			}
			refType := reflect.TypeOf(fieldValue)
			if refType == nil {
				// it can be nil if it's a pointer and it doesn't have any value!
				// so in this case we simply continue...
				continue
			}

			refTypeNative := refType
			refKind := refType.Kind()
			refVal := reflect.ValueOf(fieldValue)
			refValNative := refVal
			refIsPtr := false

			if refKind == reflect.Ptr {
				refIsPtr = true
			}

			if refIsPtr {
				refTypeNative = refType.Elem()
				if refVal.IsZero() {
					// if it's zero, we should create an empty zero value
					refValNative = reflect.Zero(refTypeNative)
				} else {
					// We should take the real indirect type value
					refValNative = reflect.Indirect(refVal)
				}
			}

			fieldValue = refValNative.Interface()
			// ============== transformations from pointer ==============\\

			if headerField.FieldName != "" && headerField.Handler == nil {

			} else if headerField.Handler != nil {
				// Check if there is a field name... and get it!
				if headerField.FieldName != "" {
					// Set the value
					exportRow.FieldValue = fieldValue
				}
				// Execute handler
				fieldValue = headerField.Handler(exportRow)
			} else {
				continue
			}

			//if fieldValue, ok := row[headerField.FieldName]; ok {
			XYAxis := headerField.XAxis + conv.IntToStr(excelRowNr)
			//log.Println("XYAxis", XYAxis, fieldValue)
			_err := f.SetCellValue(sheetName, XYAxis, fieldValue)
			if _err != nil {
				e.excelError = _err
				return false
			}
			//}
		}
	}

	// Generate a unique id
	id, _err := uuid.NewRandom()
	if _err != nil {
		e.excelError = _err
		return false
	}

	fileExtension := "xlsx"
	//fileName := id.String()

	// TODO: we should save the file id into a memory stack with id's... maybe in the filter somewhere...

	if e.ExportName == "" {
		e.ExportName = id.String()
	}

	now := time.Now()
	e.excelFileID = id
	e.excelFileName = e.ExportName + "_" + conv.Int64ToStr(now.UnixMilli())
	fullFileName := e.excelFileName + "." + fileExtension

	e.excelFullFileName = fullFileName
	e.excelFileExtension = fileExtension

	tmpPath, _err := e.GetExcelExportPath()
	if _err != nil {
		e.excelError = _err
		return false
	}

	fullFilePath := tmpPath + filesystem.DirSeparator() + fullFileName
	e.excelFullFilePath = fullFilePath

	// Set active sheet of the workbook.
	f.SetActiveSheet(sheetIndex)
	// Save spreadsheet by the given path.
	if err := f.SaveAs(fullFilePath); err != nil {
		e.excelError = _err
		return false
	}
	e.excelCreatedAt = time.Now()

	e.excelFileSizeBytes, _err = file.Size(fullFilePath)
	if _err != nil {
		e.excelError = _err
		return false
	}

	if e.SelfDeleteAfterSeconds != 0 {
		time.AfterFunc(time.Second*time.Duration(e.SelfDeleteAfterSeconds), func() {
			// Check if exists...!?
			file.Delete(fullFilePath)
		})
	}

	return true
}

func (e *Export) GetExcelFileID() uuid.UUID {
	return e.excelFileID
}

func (e *Export) GetExcelCreatedAt() time.Time {
	return e.excelCreatedAt
}

func (e *Export) GetExcelFileSizeBytes() int64 {
	return e.excelFileSizeBytes
}

func (e *Export) GetExcelFileName() string {
	return e.excelFileName
}

func (e *Export) GetExcelFullFileName() string {
	return e.excelFullFileName
}

func (e *Export) GetExcelFileExtension() string {
	return e.excelFileExtension
}

func (e *Export) GetExcelFullFilePath() string {
	return e.excelFullFilePath
}

func (e *Export) GetExcelError() error {
	return e.excelError
}

func (e *Export) QueryItems() error {
	// based on the input Filter, we should query and generate

	var _db *gorm.DB

	if e.TableName != "" {
		// TODO: should be deprecated or should be analyzed how to be used...?!...
		e.items = make(map[string]interface{})
		_db = e.Filter.DB().Table(e.TableName)
	} else if e.Model != nil {
		typeOf := reflect.TypeOf(reflect.Indirect(reflect.ValueOf(e.Model)).Interface())
		e.items = reflect.Indirect(reflect.New(reflect.SliceOf(typeOf))).Interface()
		//_db = e.Filter.DB().Model(e.Model)
		_db = e.Filter.DB()
	} else {
		_db = e.Filter.DB()
	}

	// Let's create a slice from the model

	if len(e.Preloads) > 0 {
		for _, preload := range e.Preloads {
			_db = _db.Preload(preload)
		}
	}

	//log.Println("Items2", e.Items2)
	//log.Println("Items2 len",reflect.Indirect(reflect.ValueOf(e.Items2)).Len())
	//log.Println("items", items)
	//log.Println("items len", reflect.Indirect(reflect.ValueOf(items)).Len())

	dbResult := _db.Find(&e.items)
	//log.Println("find finished...")

	if dbResult.Error != nil {
		//log.Println(dbResult.Error.Error())
		return dbResult.Error
	}

	//log.Println(e.items)
	// TODO: we should see what we have here!
	//e.nrOfRows = int64(len(e.items))
	e.nrOfRows = int64(reflect.Indirect(reflect.ValueOf(e.items)).Len())
	e.itemsSet = true
	return nil
}

func (e *Export) SetColumns(columns []ExportColumn) *Export {
	e.Columns = columns
	return e
}

func (e *Export) GetNrOfRows() int64 {
	return e.nrOfRows
}

// func (e *Export) GetItems() []map[string]interface{} {
func (e *Export) GetItems() interface{} {
	return e.items
}

// func (e *Export) SetItems(items []map[string]interface{}) *Export {
func (e *Export) SetItems(items interface{}) *Export {
	e.itemsSet = true
	e.items = items
	return e
}
