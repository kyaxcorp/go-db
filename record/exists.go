package record

import "github.com/kyaxcorp/go-helper/errors2/define"

func (r *Record) Exists() (bool, error) {
	if r.LoadData() {
		// load data procedure is ok!
		return r.isRecordExists.Get(), nil
	} else {
		// failed to load data!
		return false, define.Err(0, "load data procedure failed")
	}
}
