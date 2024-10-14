package dbresolver

import (
	"context"
	"sync"

	"github.com/kyaxcorp/go-db/driver"
	"github.com/kyaxcorp/go-helper/gor"
	"github.com/kyaxcorp/go-helper/sync/_bool"
	"github.com/kyaxcorp/go-helper/sync/_int"
	"github.com/kyaxcorp/go-logger/model"
	"gorm.io/gorm"
)

type DBResolver struct {
	*gorm.DB

	// Main config, we can read if something needed!
	mainConfig driver.Config

	configs []Config

	resolvers map[string]*resolver

	// This is the main/global/master resolver in case it doesn't match any in the resolvers...
	// global *resolver // TODO: later on it should be removed...
	// These are the main resolvers that can handle any of the operations
	masters []*resolver
	// the ones that work!
	activeMastersLock *sync.RWMutex
	activeMasters     []*resolver
	nrOfActiveMasters *_int.Int
	// the ones that don't work
	inactiveMasters     []*resolver
	nrOfInactiveMasters *_int.Int

	resolversMonitoring *gor.GInstance
	// If the monitoring is active and already scanned...
	isMonitoringActive *_bool.Bool

	ctx context.Context

	prepareStmtStore map[gorm.ConnPool]*gorm.PreparedStmtDB
	compileCallbacks []func(gorm.ConnPool) error
	// Logger
	Logger *model.Logger
}

type Config struct {
	Sources  []gorm.Dialector
	Replicas []gorm.Dialector
	Policy   Policy

	datas []interface{}
	// Logger
	Logger *model.Logger
	// Context, we use it for inside mechanisms like connection monitoring
	Ctx context.Context
}

// it's deprecated...
type failedPool struct {
	dialector gorm.Dialector
	err       error
	// This is retry market
	retry bool
}
