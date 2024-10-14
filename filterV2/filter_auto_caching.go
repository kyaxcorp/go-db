package filterV2

func (f *Input) EnableDBFieldsAutoCaching() *Input {
	f.enableDBFieldsCaching = true
	return f
}

func (f *Input) DisableDBFieldsAutoCaching() *Input {
	f.enableDBFieldsCaching = false
	return f
}
