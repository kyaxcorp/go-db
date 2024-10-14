package dbresolver

import (
	"context"
	"sync"
	"time"

	"github.com/kyaxcorp/go-helper/gor"
	"github.com/kyaxcorp/go-helper/sync/_bool"
	"github.com/kyaxcorp/go-helper/sync/_int"
	"github.com/kyaxcorp/go-helper/sync/_uint16"
	"gorm.io/gorm"
)

const ResolverAllNodesOffline = 1
const ResolverReadyToProcess = 2
const ResolverStartingUp = 3
const ResolverReadOnly = 3

type resolver struct {
	ctx context.Context

	// it means if it can handle all of the operations read/write any
	isAMaster bool

	// This is the overall status of the resolver...
	resolverStatus *_uint16.Uint16

	sources             []detailedConnPool
	nrOfSources         int
	activeSourcesLock   *sync.RWMutex
	activeSources       []detailedConnPool
	inactiveSourcesLock *sync.RWMutex
	inactiveSources     []detailedConnPool
	nrOfActiveSources   *_int.Int
	nrOfInActiveSources *_int.Int

	replicas             []detailedConnPool
	nrOfReplicas         int
	activeReplicasLock   *sync.RWMutex
	activeReplicas       []detailedConnPool
	inactiveReplicasLock *sync.RWMutex
	inactiveReplicas     []detailedConnPool
	nrOfActiveReplicas   *_int.Int
	nrOfInActiveReplicas *_int.Int
	// If the monitoring is active and already scanned...
	isMonitoringActive *_bool.Bool

	// This is the Source/Replica revival routine...
	connMonitoring *gor.GInstance

	//
	policy     Policy
	dbResolver *DBResolver
}

// This connection pool will offer much more information about the gorm.ConnPool
type detailedConnPool struct {
	// This is the connections Pool
	pool gorm.ConnPool

	// this is the latency...
	latencyNano int

	// This is the last error...
	lastErr error
	// this field shows if the pool works or not...
	isAvailable bool
	// this is the last time it was available
	lastTimeAvailable time.Time
	// this is the number of failures that it had during existence...
	nrOfFailures int
}
