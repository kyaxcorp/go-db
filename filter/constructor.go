package filter

func New(f *Input) *Input {
	if f == nil {
		f = &Input{}
	}
	defaultPageNr := DefaultPageNr
	f.PageNr = &defaultPageNr

	defaultNrOfItems := DefaultNrOfItems
	f.NrOfItems = &defaultNrOfItems

	return f
}
