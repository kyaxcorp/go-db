package record

func (r *Record) GetLastError() error {
	return r.lastError
}

func (r *Record) GetLastDBError() error {
	return r.lastDBError
}

func (r *Record) setDBError(e error) *Record {
	r.lastError = e
	r.lastDBError = e
	return r
}

func (r *Record) setError(e error) *Record {
	r.lastError = e
	return r
}
