package record

func (r *Record) callOnBeforeDelete() error {
	if _model, ok := r.modelStruct.(interface{ RecordBeforeDelete(r *Record) error }); ok {
		return _model.RecordBeforeDelete(r)
	}
	return nil
	// TODO: call other methods from record...
}

func (r *Record) callOnError() {
	if _model, ok := r.modelStruct.(interface{ RecordError(r *Record) }); ok {
		_model.RecordError(r)
	}

	// TODO: call other methods from record...
}

func (r *Record) callOnDBError() {
	if _model, ok := r.modelStruct.(interface{ RecordDBError(r *Record) }); ok {
		_model.RecordDBError(r)
	}

	// TODO: call other methods from record...
}

func (r *Record) callOnSaveError() {
	if _model, ok := r.modelStruct.(interface{ RecordSaveError(r *Record) }); ok {
		_model.RecordSaveError(r)
	}

	// TODO: call other methods from record...
}

func (r *Record) callOnDeleteError() {
	if _model, ok := r.modelStruct.(interface{ RecordDeleteError(r *Record) }); ok {
		_model.RecordDeleteError(r)
	}
	// TODO: call other methods from record...
}

func (r *Record) callOnAfterDelete() error {
	if _model, ok := r.modelStruct.(interface{ RecordAfterDelete(r *Record) error }); ok {
		return _model.RecordAfterDelete(r)
	}
	return nil
	// TODO: call other methods from record...
}

func (r *Record) callOnBeforeForceDelete() error {
	if _model, ok := r.modelStruct.(interface{ RecordBeforeForceDelete(r *Record) error }); ok {
		return _model.RecordBeforeForceDelete(r)
	}
	return nil
	// TODO: call other methods from record...
}

func (r *Record) callOnAfterForceDelete() error {
	if _model, ok := r.modelStruct.(interface{ RecordAfterForceDelete(r *Record) error }); ok {
		return _model.RecordAfterForceDelete(r)
	}
	return nil
	// TODO: call other methods from record...
}

func (r *Record) callOnBeforeSave() error {
	if _model, ok := r.modelStruct.(interface{ RecordBeforeSave(r *Record) error }); ok {
		return _model.RecordBeforeSave(r)
	}
	return nil
	// TODO: call other methods from record...
}

func (r *Record) callOnAfterDbDataLoaded() {
	if _model, ok := r.modelStruct.(interface{ RecordAfterDBDataLoaded(r *Record) }); ok {
		_model.RecordAfterDBDataLoaded(r)
	}

	// TODO: call other methods from record...
}

func (r *Record) callOnAfterSave() error {
	if _model, ok := r.modelStruct.(interface{ RecordAfterSave(r *Record) error }); ok {
		return _model.RecordAfterSave(r)
	}
	return nil
	// TODO: call other methods from record...
}

func (r *Record) callOnBeforeInsert() error {
	if _model, ok := r.modelStruct.(interface{ RecordBeforeInsert(r *Record) error }); ok {
		return _model.RecordBeforeInsert(r)
	}
	return nil
	// TODO: call other methods from record...
}

func (r *Record) callOnAfterInsert() error {
	if _model, ok := r.modelStruct.(interface{ RecordAfterInsert(r *Record) error }); ok {
		return _model.RecordAfterInsert(r)
	}
	return nil
	// TODO: call other methods from record...
}

func (r *Record) callOnBeforeUpdate() error {
	if _model, ok := r.modelStruct.(interface{ RecordBeforeUpdate(r *Record) error }); ok {
		return _model.RecordBeforeUpdate(r)
	}
	return nil
	// TODO: call other methods from record...
}

func (r *Record) callOnAfterUpdate() error {
	if _model, ok := r.modelStruct.(interface{ RecordAfterUpdate(r *Record) error }); ok {
		return _model.RecordAfterUpdate(r)
	}
	return nil
	// TODO: call other methods from record...
}
