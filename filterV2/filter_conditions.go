package filterV2

import (
	"gorm.io/gorm"
	"time"
)

func (f *Input) ApplyConditions() *Input {

	/*
		TODO: trebuie sa primim setul de field-uri! pentru care GraphQL are permisiuni/vizibilitate sa fac operatiuni!
		Chiar daca sunt mai multe cimpuri, acestea nu trebuie sa fie posibil de chemat!
	*/

	result := f.applyConditions(f.db)

	// if not nil!
	if result != nil {
		f.db = result
		// This is for counters
		// f.dbCounters = result
	}

	return f
}

func (f *Input) applyConditions(db *gorm.DB) *gorm.DB {

	/*
		TODO: trebuie sa primim setul de field-uri! pentru care GraphQL are permisiuni/vizibilitate sa fac operatiuni!
		Chiar daca sunt mai multe cimpuri, acestea nu trebuie sa fie posibil de chemat!
	*/

	return f.processGroupCondition(&GroupConditionInput{
		//ODB:            f.db,
		// We set the Original DB Client!
		DB:             db,
		GroupCondition: f.RootConditions,
	})
}

//type MultiInterface []interface{}

type QueryStatement struct {
	Query interface{}
	Args  []interface{}
}

type GroupConditionInput struct {
	// This is the original clean DB Client
	//ODB *gorm.DB
	// This is the one that modifies...
	DB             *gorm.DB
	GroupCondition *GroupCondition
}

func (f *Input) processGroupCondition(input *GroupConditionInput) *gorm.DB {
	// Check and set default values!

	// We should get the clean/original DB!
	db := input.DB
	gc := input.GroupCondition

	if gc == nil {
		return nil
	}

	// Loop through Conditions
	if gc.Conditions != nil && len(gc.Conditions) > 0 {
		for condIndex, cond := range gc.Conditions {
			// The first item OR/AND operator is ignored

			if condIndex == 0 {
			} else {
				// Add the condition
			}

			// Check all conditions, it should be only the first found!

			// TODO: we should do checks/filtration of the value and name against SQL Injections

			dbParams := &QueryStatement{}
			like := "LIKE"

			var dbFieldName string
			if cond.Equal != nil && cond.Equal.Name != "" {
				validateFieldNameAndPanic(cond.Equal.Name)
				dbFieldName = f.getDBFieldNameOrPanic(cond.Equal.Name)
				dbParams.Query = dbFieldName + " = ?"
				dbParams.Args = []interface{}{cond.Equal.Value}
				//db = db.Where(cond.Equal.Name+" = ?", cond.Equal.Value)
			} else if cond.NotEqual != nil && cond.NotEqual.Name != "" {
				validateFieldNameAndPanic(cond.NotEqual.Name)
				dbFieldName = f.getDBFieldNameOrPanic(cond.NotEqual.Name)
				dbParams.Query = dbFieldName + " != ?"
				dbParams.Args = []interface{}{cond.NotEqual.Value}
				//db = db.Where(cond.NotEqual.Name+" != ?", cond.NotEqual.Value)
			} else if cond.HigherThan != nil && cond.HigherThan.Name != "" {
				validateFieldNameAndPanic(cond.HigherThan.Name)
				dbFieldName = f.getDBFieldNameOrPanic(cond.HigherThan.Name)
				dbParams.Query = dbFieldName + " > ?"
				dbParams.Args = []interface{}{cond.HigherThan.Value}
				//db = db.Where(cond.HigherThan.Name+" > ?", cond.HigherThan.Value)
			} else if cond.HigherOrEqual != nil && cond.HigherOrEqual.Name != "" {
				validateFieldNameAndPanic(cond.HigherOrEqual.Name)
				dbFieldName = f.getDBFieldNameOrPanic(cond.HigherOrEqual.Name)
				dbParams.Query = dbFieldName + " >= ?"
				dbParams.Args = []interface{}{cond.HigherOrEqual.Value}
				//db = db.Where(cond.HigherOrEqual.Name+" >= ?", cond.HigherOrEqual.Value)
			} else if cond.LowerThan != nil && cond.LowerThan.Name != "" {
				validateFieldNameAndPanic(cond.LowerThan.Name)
				dbFieldName = f.getDBFieldNameOrPanic(cond.LowerThan.Name)
				dbParams.Query = dbFieldName + " < ?"
				dbParams.Args = []interface{}{cond.LowerThan.Value}
				//db = db.Where(cond.LowerThan.Name+" < ?", cond.LowerThan.Value)
			} else if cond.LowerOrEqual != nil && cond.LowerOrEqual.Name != "" {
				validateFieldNameAndPanic(cond.LowerOrEqual.Name)
				dbFieldName = f.getDBFieldNameOrPanic(cond.LowerOrEqual.Name)
				dbParams.Query = dbFieldName + " <= ?"
				dbParams.Args = []interface{}{cond.LowerOrEqual.Value}
				//db = db.Where(cond.LowerOrEqual.Name+" <= ?", cond.LowerOrEqual.Value)
			} else if cond.Contains != nil && cond.Contains.Name != "" {
				validateFieldNameAndPanic(cond.Contains.Name)
				dbFieldName = f.getDBFieldNameOrPanic(cond.Contains.Name)
				if cond.Contains.CaseInsensitive != nil && *cond.Contains.CaseInsensitive {
					like = "I" + like // add case-insensitive
				}
				dbParams.Query = dbFieldName + " " + like + " ?"
				dbParams.Args = []interface{}{"%" + cond.Contains.Value + "%"}
				//db = db.Where(cond.Contains.Name+" LIKE ?", "%"+cond.Contains.Value+"%")
			} else if cond.NotContains != nil && cond.NotContains.Name != "" {
				validateFieldNameAndPanic(cond.NotContains.Name)
				dbFieldName = f.getDBFieldNameOrPanic(cond.NotContains.Name)
				if cond.NotContains.CaseInsensitive != nil && *cond.NotContains.CaseInsensitive {
					like = "I" + like // add case-insensitive
				}
				dbParams.Query = dbFieldName + " NOT " + like + " ?"
				dbParams.Args = []interface{}{"%" + cond.NotContains.Value + "%"}
				//db = db.Where(cond.NotContains.Name+" NOT LIKE ?", "%"+cond.NotContains.Value+"%")
			} else if cond.BeginsWith != nil && cond.BeginsWith.Name != "" {
				validateFieldNameAndPanic(cond.BeginsWith.Name)
				dbFieldName = f.getDBFieldNameOrPanic(cond.BeginsWith.Name)
				if cond.BeginsWith.CaseInsensitive != nil && *cond.BeginsWith.CaseInsensitive {
					like = "I" + like // add case-insensitive
				}
				dbParams.Query = dbFieldName + " " + like + " ?"
				dbParams.Args = []interface{}{cond.BeginsWith.Value + "%"}
				//db = db.Where(cond.BeginsWith.Name+" LIKE ?", cond.BeginsWith.Value+"%")
			} else if cond.NotBeginsWith != nil && cond.NotBeginsWith.Name != "" {
				validateFieldNameAndPanic(cond.NotBeginsWith.Name)
				dbFieldName = f.getDBFieldNameOrPanic(cond.NotBeginsWith.Name)
				if cond.NotBeginsWith.CaseInsensitive != nil && *cond.NotBeginsWith.CaseInsensitive {
					like = "I" + like // add case-insensitive
				}
				dbParams.Query = dbFieldName + " NOT " + like + " ?"
				dbParams.Args = []interface{}{cond.NotBeginsWith.Value + "%"}
				//db = db.Where(cond.BeginsWith.Name+" LIKE ?", cond.BeginsWith.Value+"%")
			} else if cond.EndsWith != nil && cond.EndsWith.Name != "" {
				validateFieldNameAndPanic(cond.EndsWith.Name)
				dbFieldName = f.getDBFieldNameOrPanic(cond.EndsWith.Name)
				if cond.EndsWith.CaseInsensitive != nil && *cond.EndsWith.CaseInsensitive {
					like = "I" + like // add case-insensitive
				}
				dbParams.Query = dbFieldName + " " + like + " ?"
				dbParams.Args = []interface{}{"%" + cond.EndsWith.Value}
				//db = db.Where(cond.EndsWith.Name+" LIKE ?", "%"+cond.EndsWith.Value)
			} else if cond.NotEndsWith != nil && cond.NotEndsWith.Name != "" {
				validateFieldNameAndPanic(cond.NotEndsWith.Name)
				dbFieldName = f.getDBFieldNameOrPanic(cond.NotEndsWith.Name)
				if cond.NotEndsWith.CaseInsensitive != nil && *cond.NotEndsWith.CaseInsensitive {
					like = "I" + like // add case-insensitive
				}
				dbParams.Query = dbFieldName + " NOT " + like + " ?"
				dbParams.Args = []interface{}{"%" + cond.NotEndsWith.Value}
				//db = db.Where(cond.NotEndsWith.Name+" NOT LIKE ?", "%"+cond.NotEndsWith.Value)
			} else if cond.In != nil && cond.In.Name != "" {
				validateFieldNameAndPanic(cond.In.Name)
				dbFieldName = f.getDBFieldNameOrPanic(cond.In.Name)
				dbParams.Query = dbFieldName + " IN (?)"
				dbParams.Args = []interface{}{cond.In.Value}
				//db = db.Where(cond.In.Name+" IN (?)", cond.In.Value)
			} else if cond.NotIn != nil && cond.NotIn.Name != "" {
				validateFieldNameAndPanic(cond.NotIn.Name)
				dbFieldName = f.getDBFieldNameOrPanic(cond.NotIn.Name)
				dbParams.Query = dbFieldName + " NOT IN (?)"
				dbParams.Args = []interface{}{cond.NotIn.Value}
				//db = db.Where(cond.NotIn.Name+" NOT IN (?)", cond.NotIn.Value)
			} else if cond.IsNull != nil && cond.IsNull.Name != "" {
				validateFieldNameAndPanic(cond.IsNull.Name)
				dbFieldName = f.getDBFieldNameOrPanic(cond.IsNull.Name)
				dbParams.Query = dbFieldName + " IS NULL"
				dbParams.Args = []interface{}{}
				//db = db.Where(cond.IsNull.Name + " IS NULL")
			} else if cond.IsNotNull != nil && cond.IsNotNull.Name != "" {
				validateFieldNameAndPanic(cond.IsNotNull.Name)
				dbFieldName = f.getDBFieldNameOrPanic(cond.IsNotNull.Name)
				dbParams.Query = dbFieldName + " NOT NULL"
				dbParams.Args = []interface{}{}
				//db = db.Where(cond.IsNotNull.Name + " NOT NULL")
			} else if cond.IsEmpty != nil && cond.IsEmpty.Name != "" {
				validateFieldNameAndPanic(cond.IsEmpty.Name)
				dbFieldName = f.getDBFieldNameOrPanic(cond.IsEmpty.Name)
				dbParams.Query = "(" + dbFieldName + " ='' OR " + dbFieldName + " IS NULL)"
				dbParams.Args = []interface{}{}
				//db = db.Where("(" + cond.IsEmpty.Name + " ='' OR " + cond.IsEmpty.Name + " IS NULL)")
			} else if cond.IsNotEmpty != nil && cond.IsNotEmpty.Name != "" {
				validateFieldNameAndPanic(cond.IsNotEmpty.Name)
				dbFieldName = f.getDBFieldNameOrPanic(cond.IsNotEmpty.Name)
				dbParams.Query = "(" + dbFieldName + " !='' AND " + dbFieldName + " NOT NULL)"
				dbParams.Args = []interface{}{}
				//db = db.Where("(" + cond.IsNotEmpty.Name + " !='' AND " + cond.IsNotEmpty.Name + " NOT NULL)")
			} else if cond.IsTrue != nil && cond.IsTrue.Name != "" {
				validateFieldNameAndPanic(cond.IsTrue.Name)
				dbFieldName = f.getDBFieldNameOrPanic(cond.IsTrue.Name)
				dbParams.Query = dbFieldName + " IS TRUE"
				dbParams.Args = []interface{}{}
				//db = db.Where(cond.IsTrue.Name + " IS TRUE")
			} else if cond.IsFalse != nil && cond.IsFalse.Name != "" {
				validateFieldNameAndPanic(cond.IsFalse.Name)
				dbFieldName = f.getDBFieldNameOrPanic(cond.IsFalse.Name)
				dbParams.Query = dbFieldName + " IS NOT TRUE"
				dbParams.Args = []interface{}{}
				//db = db.Where(cond.IsFalse.Name + " IS NOT TRUE")
			} else if cond.Between != nil && cond.Between.Name != "" {
				validateFieldNameAndPanic(cond.Between.Name)
				dbFieldName = f.getDBFieldNameOrPanic(cond.Between.Name)
				dbParams.Query = "(" + dbFieldName + " BETWEEN ? AND ? )"
				dbParams.Args = []interface{}{cond.Between.Start, cond.Between.End}
				//db = db.Where("("+cond.Between.Name+" BETWEEN ? AND ? )", cond.Between.Start, cond.Between.End)
			} else if cond.NotBetween != nil && cond.NotBetween.Name != "" {
				validateFieldNameAndPanic(cond.NotBetween.Name)
				dbFieldName = f.getDBFieldNameOrPanic(cond.NotBetween.Name)
				dbParams.Query = "(" + dbFieldName + " NOT BETWEEN ? AND ? )"
				dbParams.Args = []interface{}{cond.NotBetween.Start, cond.NotBetween.End}
				//db = db.Where("("+cond.NotBetween.Name+" NOT BETWEEN ? AND ? )", cond.NotBetween.Start, cond.NotBetween.End)
			} else if cond.BetweenUnixTimestamp != nil && cond.BetweenUnixTimestamp.Name != "" {
				validateFieldNameAndPanic(cond.BetweenUnixTimestamp.Name)
				dbFieldName = f.getDBFieldNameOrPanic(cond.BetweenUnixTimestamp.Name)
				dbParams.Query = "(" + dbFieldName + " BETWEEN ? AND ? )"
				dbParams.Args = []interface{}{
					time.Unix(*cond.BetweenUnixTimestamp.Start, 0).Format("2006-01-02 15:04:05"),
					time.Unix(*cond.BetweenUnixTimestamp.End, 0).Format("2006-01-02 15:04:05"),
				}
				//db = db.Where("("+cond.Between.Name+" BETWEEN ? AND ? )", cond.Between.Start, cond.Between.End)
			} else if cond.NotBetweenUnixTimestamp != nil && cond.NotBetweenUnixTimestamp.Name != "" {
				validateFieldNameAndPanic(cond.NotBetweenUnixTimestamp.Name)
				dbFieldName = f.getDBFieldNameOrPanic(cond.NotBetweenUnixTimestamp.Name)
				dbParams.Query = "(" + dbFieldName + " NOT BETWEEN ? AND ? )"
				dbParams.Args = []interface{}{
					time.Unix(*cond.NotBetweenUnixTimestamp.Start, 0).Format("2006-01-02 15:04:05"),
					time.Unix(*cond.NotBetweenUnixTimestamp.End, 0).Format("2006-01-02 15:04:05"),
				}
				//db = db.Where("("+cond.NotBetween.Name+" NOT BETWEEN ? AND ? )", cond.NotBetween.Start, cond.NotBetween.End)
			}

			if cond.Or != nil && *cond.Or && condIndex != 0 {
				// If it's OR operator
				db = db.Or(dbParams.Query, dbParams.Args...)
			} else {
				// it's AND Operator
				db = db.Where(dbParams.Query, dbParams.Args...)
			}

			// TODO: do also for JSON values and types!
		}
	}

	// Loop through Groups
	if gc.Groups != nil && len(gc.Groups) > 0 {
		for groupIndex, group := range gc.Groups {
			// The first item OR/AND operator is ignored

			// Check the group
			groupDb := f.processGroupCondition(&GroupConditionInput{
				// here we should give a clean DB!
				DB:             input.DB,
				GroupCondition: group,
			})

			if groupDb == nil {
				continue
			}

			// Group conditions
			if group.Or != nil && *group.Or && groupIndex != 0 {
				// If it's OR operator
				// TODO: is it working as expected?! i think no... because it's not always adding parentheses
				db = db.Or(groupDb)
			} else {
				// it's AND Operator
				db = db.Where(groupDb)
			}
		}
	}

	return db
}
