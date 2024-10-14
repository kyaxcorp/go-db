package dbinstance

import (
	"gorm.io/gorm"
	"sync"
)

type OnMissingAutoCreate func(instanceName string) (*gorm.DB, error)

type Instance struct {
	// We can't use here because if multiple goroutines are accessing same lock with RLock, it will not save
	// us from creating a new instance, we can have a panic error of nil point dereference
	//lock      sync.RWMutex
	lock      sync.Mutex
	instances map[string]*gorm.DB
	// Define the callback in case of missing instance
	OnMissingAutoCreate OnMissingAutoCreate
}
