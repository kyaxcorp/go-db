package dbresolver

import "github.com/kyaxcorp/go-helper/_context"

func (r *resolver) getActiveSources() []detailedConnPool {
	r.activeSourcesLock.RLock()
	defer r.activeSourcesLock.RUnlock()
	return r.activeSources
}

func (r *resolver) getInactiveSources() []detailedConnPool {
	r.inactiveSourcesLock.RLock()
	defer r.inactiveSourcesLock.RUnlock()
	return r.inactiveSources
}

func (r *resolver) getActiveReplicas() []detailedConnPool {
	r.activeReplicasLock.RLock()
	defer r.activeReplicasLock.RUnlock()
	return r.activeReplicas
}

func (r *resolver) getInactiveReplicas() []detailedConnPool {
	r.inactiveReplicasLock.RLock()
	defer r.inactiveReplicasLock.RUnlock()
	return r.inactiveReplicas
}

func (r *resolver) getNrOfSources() int {
	return r.nrOfSources
}

func (r *resolver) getNrOfReplicas() int {
	return r.nrOfReplicas
}

func (r *resolver) IsTerminating() bool {
	// We simply create a temporary context to not be blocked by using directly
	tmpContext := _context.WithCancel(r.ctx)
	if tmpContext.IsDone() {
		return true
	}
	return false
}
