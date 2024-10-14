package dbresolver

import "github.com/kyaxcorp/go-db/driver"

func (dr *DBResolver) SetMainConfig(config driver.Config) {
	dr.mainConfig = config
}
