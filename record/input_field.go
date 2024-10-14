package record

// GetInputFieldValue -> will return the value that has been set!
func (r *Record) GetInputFieldValue(fieldName string) (interface{}, error) {
	//switch r.inputDataType {
	//case inputDataMapInterface:
	//	if val, ok := r.dataMap[fieldName]; ok {
	//		return val, nil
	//	}
	//	return nil, define.Err(0, "field doesn't exist ->", fieldName)
	//case inputDataStruct:
	//	//if !r.dataStrHelper.FieldExists(fieldName) {
	//	//	return nil, define.Err(0, "field doesn't exist ->", fieldName)
	//	//}
	//	return r.dataStrHelper.GetFieldValue(fieldName), nil
	//}

	return r.dataStrHelper.GetFieldValue(fieldName), nil
	//return nil, nil
}
