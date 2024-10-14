package filter

import (
	"strconv"

	"github.com/kyaxcorp/go-db/helper"
	"github.com/kyaxcorp/go-helper/_struct"
	"github.com/kyaxcorp/go-helper/errors2/define"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	//"log"
	"strings"
	"sync"
)

type InputModel struct {
	Name  string
	Model interface{}

	// TableName -> used when using Joins, which automatically makes a relation to a custom table name!
	TableName string
}

func (f *Input) Apply() *Input {
	f.check()
	f.ApplyConditions()
	f.applyDefaultScope()
	f.ApplyPagination()
	f.ApplyOrdering()
	return f
}

func (f *Input) check() {
	// Check and set default values!
	if f.PageNr == nil {
		f.PageNr = new(int64)
		*f.PageNr = DefaultPageNr
	} else if *f.PageNr <= 0 {
		*f.PageNr = DefaultPageNr
	}

	if f.NrOfItems == nil {
		f.NrOfItems = new(int64)
		*f.NrOfItems = DefaultNrOfItems
		//} else if *f.NrOfItems <= 0 {
	} else if *f.NrOfItems == -1 {
		// Don't do anything... -1 it's unlimited
	} else if *f.NrOfItems <= 0 {
		*f.NrOfItems = DefaultNrOfItems
	}

	if f.maxNrOfItems == nil {
		f.maxNrOfItems = new(int64)
		*f.maxNrOfItems = DefaultMaxNrOfItems
	} else if *f.maxNrOfItems == -1 {
		// Don't do anything... -1 it's unlimited
	} else if *f.maxNrOfItems <= 0 {
		*f.maxNrOfItems = DefaultMaxNrOfItems
	}

	// Check if Max allows unlimited!
	if *f.NrOfItems == -1 && *f.maxNrOfItems != -1 {
		// reset the nr of items to default value because it's not allowed
		*f.NrOfItems = DefaultNrOfItems
	}

	f.checkContext()
}

/*type InputModel interface {
	// Each GORM NonPtrObj has this function or should have at least...
	TableName() string
}*/

func (f *Input) SetModels(models ...InputModel) *Input {
	if f.models == nil {
		f.models = make(map[string]cachedModel)
	}
	if f.cachedDBFields == nil {
		f.cachedDBFields = make(map[string]DBField)
	}

	//f.db
	//f.db.NewScope(model).TableName()

	/*
		TODO:
			We should represent:
			- What's the model name in GraphQL (as the user sees it)
			- What's the model structure to which is bind!
			-
	*/

	// Parse the model table name

	modelNr := 0
	for _, inputModel := range models {
		modelNr++
		// Get the model native Type/Name

		model := inputModel.Model
		modelName := _struct.GetName(model)

		// Get the model DB Table Name
		tableName, _err := helper.GetModelDBTableName(f.db, model)
		if _err != nil {
			panic("failed to parse model table name")
		}
		// Save the model as it is

		if inputModel.TableName != "" {
			tableName = inputModel.TableName
		}

		inputModelName := strings.ToLower(inputModel.Name)

		if f.primaryModelName == "" {
			f.primaryModelName = inputModelName
		}

		modelCache := cachedModel{
			modelName:   modelName,
			model:       model,
			dbColumns:   make(map[string]string),
			dbTableName: tableName,
		}

		f.models[inputModelName] = modelCache

		if modelNr == 1 {
			f.primaryModel = modelCache
		}

		// Let's also parse model Columns!
		//modelFields := _struct.GetExportableFieldsByJSONTag(model)
		//lowerCaseModelFields := make(map)
		//for _,fieldName := range modelFields{
		//
		//}

		s, err := schema.Parse(model, &sync.Map{}, schema.NamingStrategy{})
		if err != nil {
			panic("failed to create schema")
		}

		//m := make(map[string]string)
		found := make(map[string]bool)
		for _, field := range s.Fields {
			dbFieldName := field.DBName
			quotedDbFieldName := strconv.Quote(field.DBName)
			modelFieldName := field.Name
			//m[modelFieldName] = dbFieldName

			filteredDbFieldName := strings.ToLower(modelFieldName)

			// We will lowercase the field name from the structure/model!
			f.models[inputModelName].dbColumns[filteredDbFieldName] = dbFieldName

			tableNameFieldName := tableName + "." + dbFieldName
			quotedTableNameFieldName := strconv.Quote(tableName) + "." + strconv.Quote(dbFieldName)

			dbField := DBField{
				FieldName:                dbFieldName,
				QuotedFieldName:          quotedDbFieldName,
				TableNameFieldName:       tableNameFieldName,
				QuotedTableNameFieldName: quotedTableNameFieldName,
			}

			if f.enableDBFieldsCaching {
				// Save as it is
				// Here, we can have multiple models, so, we should set only the first ones that we find!
				// Multiple models can have same fields, so let's do a check if the field already exists!
				if _, ok := f.cachedDBFields[filteredDbFieldName]; !ok {
					f.cachedDBFields[filteredDbFieldName] = dbField
				}
				// Save with model name and field name
				f.cachedDBFields[inputModelName+"."+filteredDbFieldName] = dbField
			}

			found[modelFieldName] = true

			//log.Println(modelFieldName + " -> " + dbFieldName)
		}

		// Try also getting from tags...

		/*mHelper := _struct.New(model)
		fields := mHelper.Map()
		for fieldName, _ := range fields {
			if _, ok := found[fieldName]; !ok {
				dbFieldName := mHelper.GetFieldTagKeyValue(fieldName, "gorm", "column")
				if dbFieldName != "" {
					// it exists, add manually

					modelFieldName := fieldName
					filteredDbFieldName := strings.ToLower(modelFieldName)

					// We will lowercase the field name from the structure/model!
					f.models[inputModelName].dbColumns[filteredDbFieldName] = dbFieldName

					if f.enableDBFieldsCaching {
						// Save as it is
						f.cachedDBFields[filteredDbFieldName] = dbFieldName
						// Save with model name and field name
						f.cachedDBFields[inputModelName+"."+filteredDbFieldName] = dbFieldName
					}

					found[modelFieldName] = true
				}
			}

		}*/

		//

		//
		/*result, _ := f.db.Migrator().ColumnTypes(model)
		for _, v := range result {
			// TODO: what to do with the column tag which can change entirely the name
			// 		for which we can't reverse the name...
			//		we need the function column types to return the original field name for which the conversion
			//		has been made!
			// 		or we should loop by ourselves and get the names

			dbFieldName := v.Name()
			filteredDbFieldName := strings.ReplaceAll(dbFieldName, "_", "")
			//log.Println(v.Name(), v.DatabaseTypeName(), filteredDbFieldName)

			f.models[inputModelName].dbColumns[filteredDbFieldName] = dbFieldName
			// Let's get it back to original name and after that map!
			//log.Println()
		}*/
	}

	return f
}

func (f *Input) getDBFieldName(fieldName string) (string, error) {
	// transform
	/*
		1. Check if there is a DOT (.) in the string
		if there is one, check the by this model name!
		Loop there to find the field name! if there is no
		field, then throw an error that field doesn't exist!

		2. If there is no DOT (.) search consecutively through
		all existing models for this field, the first found, first delivered...
		3. if still not found then throw an error...

	*/

	lowerFieldName := strings.ToLower(fieldName)

	//  Check cache here...
	if dbField, ok := f.cachedDBFields[lowerFieldName]; ok {
		//return dbField.TableNameFieldName, nil
		return dbField.QuotedTableNameFieldName, nil
	}

	if strings.Contains(fieldName, ".") {
		// Case 1
		splitted := strings.Split(lowerFieldName, ".")
		modelName := splitted[0]
		fName := splitted[1]
		if _, ok := f.models[modelName]; !ok {
			return "", define.Err(1, "this model doesn't exist -> "+modelName)
		}
		if _, ok := f.models[modelName].dbColumns[fName]; !ok {
			return "", define.Err(2, "this field doesn't exist -> "+fName)
		}

		//return f.models[modelName].dbColumns[fName], nil
		//
		// 20.05.2022 - am modificat pentru cashpot la stations (acesta are join cu inventory)
		// dadea eroare in filtre pentru coloana location_id -> is ambiguous
		// si aici ca solutie e sa adaug table name inainte de field, pentru ca sa concretizez acesta....
		//
		// 30.06.2022 -> am adaugat Double Quoting... "" pentru ca cind faci Joins("CreatedBy") si faci Search CreatedBy.first_name
		// dar eroare : ERROR: no data source matches prefix: updatedby in this context (SQLSTATE 42P01)
		// Gorm by default adauga doudble quotes

		return strconv.Quote(f.models[modelName].dbTableName) + "." + strconv.Quote(f.models[modelName].dbColumns[fName]), nil
	} else {
		// Case 2

		found := false
		dbFieldName := ""
		dbTableName := ""
		for _, model := range f.models {
			//log.Println("model", model.dbColumns)
			if dbField, ok := model.dbColumns[lowerFieldName]; ok {
				found = true
				dbFieldName = dbField
				dbTableName = model.dbTableName
				break
			}
		}
		if !found {
			return "", define.Err(3, "this field doesn't exist -> "+lowerFieldName)
		}
		// 30.06.2022 -> am adaugat Double Quoting... "" pentru ca cind faci Joins("CreatedBy") si faci Search CreatedBy.first_name
		// dar eroare : ERROR: no data source matches prefix: updatedby in this context (SQLSTATE 42P01)
		// Gorm by default adauga doudble quotes
		return strconv.Quote(dbTableName) + "." + strconv.Quote(dbFieldName), nil
	}
}

func (f *Input) getDBFieldNameOrPanic(fieldName string) string {
	fName, _err := f.getDBFieldName(fieldName)
	if _err != nil {
		panic(_err.Error())
	}
	return fName
}

func (f *Input) DB() *gorm.DB {
	return f.db
}
