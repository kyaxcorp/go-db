package dbinstance

import (
	"gorm.io/gorm"
	"sync"
)

func NewInstance() *Instance {
	return &Instance{
		//lock:      sync.RWMutex{},
		lock:      sync.Mutex{},
		instances: make(map[string]*gorm.DB),
	}
}
