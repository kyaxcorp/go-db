package dbresolver

import (
	"sync"

	"github.com/kyaxcorp/go-db/codes"
	"github.com/kyaxcorp/go-helper/errors2/define"
	"github.com/kyaxcorp/go-helper/sync/_bool"
	"github.com/kyaxcorp/go-helper/sync/_int"
	"github.com/kyaxcorp/go-helper/sync/_uint16"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func (dr *DBResolver) compile() error {
	// we can have multiple resolvers (with different configurations)
	// each resolver can have its own purpose...

	// Start a Resolvers Monitoring system... which will switch which of the resolvers will be the primary one...

	_error := func() *zerolog.Event {
		return dr.LErrorF("dbr.compile")
	}

	for _, config := range dr.configs {
		if err := dr.compileConfig(config); err != nil {
			return err
		}
	}

	// we are starting resolvers monitoring later after the resolvers are started...
	// why, because the resolvers monitoring requires statuses of the resolvers
	// so the first 2 seconds or more you may see a delay before receiving the response...
	_err := dr.startResolversMonitoring()
	if _err != nil {
		_error().Err(codes.ErrFailedToStartResolversMonitoring).Msg("")
		return codes.ErrFailedToStartResolversMonitoring
	}

	return nil
}

// compileConfig - compiling resolver, each resolver is compiled independently
// each resolver will have its own monitoring service
func (dr *DBResolver) compileConfig(config Config) (err error) {
	// this is the primary connPool
	var connPool = dr.DB.Config.ConnPool

	info := func() *zerolog.Event {
		return dr.LInfoF("dbr.compileConfig")
	}
	/*warn := func() *zerolog.Event {
		return dr.LWarnF("dbr.compileConfig")
	}*/
	_error := func() *zerolog.Event {
		return dr.LErrorF("dbr.compileConfig")
	}

	info().Msg("compiling configuration")

	// Creating the resolver...
	var r = resolver{
		ctx:        config.Ctx,
		dbResolver: dr,
		policy:     config.Policy,
		// Context is being needed for wait functions that boolean has...
		isMonitoringActive: _bool.NewValContext(false, config.Ctx),

		resolverStatus: _uint16.NewVal(ResolverStartingUp),

		nrOfActiveSources:    _int.NewVal(0),
		nrOfInActiveSources:  _int.NewVal(0),
		nrOfActiveReplicas:   _int.NewVal(0),
		nrOfInActiveReplicas: _int.NewVal(0),

		activeSourcesLock:    &sync.RWMutex{},
		inactiveSourcesLock:  &sync.RWMutex{},
		activeReplicasLock:   &sync.RWMutex{},
		inactiveReplicasLock: &sync.RWMutex{},
	}

	if preparedStmtDB, ok := connPool.(*gorm.PreparedStmtDB); ok {
		connPool = preparedStmtDB.ConnPool
	}

	if len(config.Sources) == 0 {
		//r.sources = []gorm.ConnPool{connPool}
		// if there are no defined sources, we will set the primary one
		// as the default source
		r.sources = []detailedConnPool{{pool: connPool}}
	} else {
		detailedConnPools := dr.convertToDetailedConnPool(config.Sources)

		if len(detailedConnPools) == 0 {
			// Error will be returned only when it's critical!
			return define.Err(0, "no pools are available")
		}

		r.sources = detailedConnPools
	}

	// Count & set how many sources...
	r.nrOfSources = len(r.sources)
	r.replicas = dr.convertToDetailedConnPool(config.Replicas)
	// Count & set how many replicas...
	r.nrOfReplicas = len(r.replicas)

	// Start the connection monitoring
	_err := r.startConnectionsMonitoring()
	if _err != nil {
		_error().Err(codes.ErrFailedToStartConnectionsMonitoringForResolver).Msg("")
		return codes.ErrFailedToStartConnectionsMonitoringForResolver
	}

	// check if there are any defined Tables in the configuration
	// it's the tables names or structures
	// if there are, then we should save as MAP table->resolver (the defined table, and the current created resolver)
	if len(config.datas) > 0 {
		for _, data := range config.datas {
			if t, ok := data.(string); ok {
				dr.resolvers[t] = &r
			} else {
				stmt := &gorm.Statement{DB: dr.DB}
				if err := stmt.Parse(data); err == nil {
					dr.resolvers[stmt.Table] = &r
				} else {
					return err
				}
			}
		}
	} else {
		// set to masters
		r.isAMaster = true
		dr.masters = append(dr.masters, &r)

		// set the global one...maybe later on it will be deprecated... and we will use the masters
		// dr.global = &r // TODO: should we remove this?
	}

	for _, fc := range dr.compileCallbacks {
		if err = r.call(fc); err != nil {
			return err
		}
	}

	return nil
}
