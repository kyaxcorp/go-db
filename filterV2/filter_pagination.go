package filterV2

import (
	"github.com/kyaxcorp/go-db/scope"
	"gorm.io/gorm"
)

func (f *Input) ApplyPagination() *Input {
	f.db = f.db.Scopes(f.applyPagination)
	return f
}

func (f *Input) applyPagination(db *gorm.DB) *gorm.DB {
	// let's check if max Nr Of Items is unlimited!
	if *f.NrOfItems == -1 && *f.maxNrOfItems == -1 {
		// don't do anything...
		// we allow to query everything at once...
		return db
	}

	db = db.Scopes(
		scope.Paginate(int(*f.PageNr), int(*f.NrOfItems), int(*f.maxNrOfItems)),
	)
	return db
}

// SetNrOfItems -> this is the requested nr of items!
func (f *Input) SetNrOfItems(nrOfItems int64) *Input {
	f.NrOfItems = &nrOfItems
	return f
}

// SetMaxNrOfItems -> this is the max allowed value that can be set when requesting... or
// the max value that will be taken in case of if NrOfItems is higher than this limit!
// -1 can also be set! -> this means it's ALL or unlimited!
func (f *Input) SetMaxNrOfItems(maxNrOfItems int64) *Input {
	f.maxNrOfItems = &maxNrOfItems
	return f
}

func (f *Input) SetPageNr(pageNr int64) *Input {
	f.PageNr = &pageNr
	return f
}
