package filterV2

import (
	"gorm.io/gorm"
)

func (f *Input) SetDB(db *gorm.DB) *Input {
	//f.db = db
	return f.SetDb(db)
}

// Because when calling in each method from gorm, they call a function called
// get instance, and this function is cloning it! that's why we need to call back the get db client from
// the Filter!
// a new address for the pointer it's been set each time! the old address is not touched anymore!

func (f *Input) SetDb(dbClient *gorm.DB) *Input {
	f.SetMainDB(dbClient)
	f.SetCounterDB(dbClient)
	return f
}

func (f *Input) SetMainDB(dbClient *gorm.DB) *Input {
	// Clone the DB Gorm instance
	// And also clone the Statement!
	f.db = dbClient.WithContext(f.ctx)

	return f
}

func (f *Input) SetCounterDB(dbClient *gorm.DB) *Input {
	// Clone the DB Gorm instance
	// And also clone the Statement!
	//f.db = dbClient.Attrs()
	f.dbCounters = dbClient.WithContext(f.ctx)
	return f
}
