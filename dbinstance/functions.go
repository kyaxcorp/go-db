package dbinstance

import (
	"gorm.io/gorm"
)

func (i *Instance) SaveClientToInstances(
	instanceName string,
	client *gorm.DB,
) {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.instances[instanceName] = client
}

// GetClientByInstanceId - returns an existing/saved instance
func (i *Instance) GetClientByInstanceId(
	instanceName string,
) (*gorm.DB, error) {
	// Check from cache first!
	//i.lock.RLock()
	//defer i.lock.RUnlock()
	i.lock.Lock()
	defer i.lock.Unlock()
	if client, ok := i.instances[instanceName]; ok {
		return client, nil
	} else {
		client, _err := i.OnMissingAutoCreate(instanceName)
		// if failed to create client, return the error...
		if _err != nil {
			return nil, _err
		}
		// Save the client...
		i.instances[instanceName] = client
		return client, nil
	}
}
