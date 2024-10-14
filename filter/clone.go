package filter

func (f *Input) NewInstanceWithConditions() *Input {
	pageNr := *f.PageNr
	nrOfItems := *f.NrOfItems
	maxNrOfItems := *f.maxNrOfItems

	// we don't clone entirely the instance...

	newInput := &Input{
		PageNr:       &pageNr,
		NrOfItems:    &nrOfItems,
		maxNrOfItems: &maxNrOfItems,

		// we should remove the pointer from these vars and then create new ones
		// temporarily will work this, because we don't change the conditions in the process...
		// TODO: clone as it should!
		Order:          f.Order,
		Search:         f.Search,
		RootConditions: f.RootConditions,
	}
	return newInput
}
