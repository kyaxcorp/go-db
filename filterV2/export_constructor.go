package filterV2

func NewExport(e *Export) *Export {
	if e == nil {
		e = &Export{}
	}

	return e
}
