package dbresolver

import (
	"sync"
	"time"

	"gorm.io/gorm"
)

// It will return the ones that are ok and the ones that received error on connection...
// the ones that receive error on connection will be handled by a special reconnector which will always
// retry until success and add back to the pool!
func (dr *DBResolver) convertToConnPoolAvailable(
	dialectors []gorm.Dialector,
) (connPools []gorm.ConnPool, failedPools []failedPool) {
	config := *dr.DB.Config
	for _, dialector := range dialectors {
		// Connect to the DB, if it has not connected to, then the connection pool will not be created!
		// If it hasn't connected to one of the DB, then it will return an error...
		// TODO: we should convert all possible, and the ones that failed we should return back!
		if db, err := gorm.Open(dialector, &config); err == nil {
			connPool := db.Config.ConnPool
			if preparedStmtDB, ok := connPool.(*gorm.PreparedStmtDB); ok {
				connPool = preparedStmtDB.ConnPool
			}

			dr.prepareStmtStore[connPool] = &gorm.PreparedStmtDB{
				ConnPool: db.Config.ConnPool,
				Stmts:    map[string]*gorm.Stmt{},
				Mux:      &sync.RWMutex{},
				//PreparedSQL: make([]string, 0, 100),
			}

			// Gather the connected pools
			connPools = append(connPools, connPool)
		} else {
			// TODO: detect what type of error is this:
			// in case it's invalid credentials, then we should mark as invalid and it shouldn't try anymore reconnecting
			// if it couldn't connect, then it should remain as failed Pool but with a retry marker
			// to detect if there are invalid credentials, we should get into the Each Driver separately...
			// each driver has it's own error codes...
			// so we will leave each error to be treated as retry

			// Gather the failed pools
			failedPools = append(failedPools, failedPool{
				dialector: dialector,
				err:       err,
				retry:     true,
			})
		}
	}

	return connPools, failedPools
}

func (dr *DBResolver) convertToDetailedConnPool(
	dialectors []gorm.Dialector,
) (connPools []detailedConnPool) {
	config := *dr.DB.Config
	for _, dialector := range dialectors {
		// Connect to the DB, if it has not connected to, then the connection pool will not be created!
		// If it hasn't connected to one of the DB, then it will return an error...
		// TODO: we should convert all possible, and the ones that failed we should return back!
		isAvailable := false

		db, err := gorm.Open(dialector, &config)
		connPool := db.Config.ConnPool
		if err == nil {
			isAvailable = true
		} else {
			isAvailable = false
		}

		if preparedStmtDB, ok := connPool.(*gorm.PreparedStmtDB); ok {
			connPool = preparedStmtDB.ConnPool
		}

		dr.prepareStmtStore[connPool] = &gorm.PreparedStmtDB{
			ConnPool: db.Config.ConnPool,
			Stmts:    map[string]*gorm.Stmt{},
			Mux:      &sync.RWMutex{},
			//PreparedSQL: make([]string, 0, 100),
		}

		connPools = append(connPools, detailedConnPool{
			pool:              connPool,
			isAvailable:       isAvailable,
			lastTimeAvailable: time.Now(),
		})
	}

	return connPools
}

func (dr *DBResolver) convertToConnPool(dialectors []gorm.Dialector) (connPools []gorm.ConnPool, err error) {
	config := *dr.DB.Config
	for _, dialector := range dialectors {
		// Connect to the DB, if it has not connected to, then the connection pool will not be created!
		// If it hasn't connected to one of the DB, then it will return an error...
		// TODO: we should convert all possible, and the ones that failed we should return back!
		if db, err := gorm.Open(dialector, &config); err == nil {
			connPool := db.Config.ConnPool
			if preparedStmtDB, ok := connPool.(*gorm.PreparedStmtDB); ok {
				connPool = preparedStmtDB.ConnPool
			}

			dr.prepareStmtStore[connPool] = &gorm.PreparedStmtDB{
				ConnPool: db.Config.ConnPool,
				Stmts:    map[string]*gorm.Stmt{},
				Mux:      &sync.RWMutex{},
				//PreparedSQL: make([]string, 0, 100),
			}

			connPools = append(connPools, connPool)
		} else {
			return nil, err
		}
	}

	return connPools, err
}
