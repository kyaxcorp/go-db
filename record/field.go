package record

func (r *Record) GetDBFieldName(fieldName string) string {
	// TODO: should we check for existence?
	return r.modelFieldNamings[fieldName]
}

func (r *Record) IsStructFieldExists(fieldName string) bool {
	if _, ok := r.modelFieldNamings[fieldName]; ok {
		return true
	}
	return false
}
