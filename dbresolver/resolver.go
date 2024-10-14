package dbresolver

import (
	"time"

	"github.com/gookit/color"
	"github.com/kyaxcorp/go-db/codes"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func (r *resolver) resolve(
	stmt *gorm.Statement,
	op Operation,
) (connPool gorm.ConnPool, err error) {
	funcStartTime := time.Now().Nanosecond()
	info := func() *zerolog.Event {
		return r.LInfoF("r.resolve")
	}
	defer printConsumedTime(info, funcStartTime)

	_error := func() *zerolog.Event {
		return r.LErrorF("r.resolve")
	}

	info().Msg("resolving...")
	defer info().Msg("leaving...")

	info().Msg("waiting for monitoring to get online...")
	if r.isMonitoringActive.WaitUntilTrue().Get() == false {
		// maybe it's been canceled!
		if r.IsTerminating() {
			return nil, codes.ErrTerminating
		}
	}

	info().Msg("generating resolve options")
	_resolveOptions := &resolveOptions{
		context:    stmt.Context,
		dbResolver: r.dbResolver,
		resolver:   r,
	}

	// If the operation is read and there are any replicas, then try
	maxRetries := 3
	retryDelayMs := 1000
	currentNrOfRetries := 0

	retryConnection := func(caseNr int) bool {
		_error().
			Int("case_nr", caseNr). // It shows where it captured this event
			Int("retry_nr", currentNrOfRetries).
			Int("max_retries", maxRetries).
			Msg("failed to find an active/restored connection pool")

		if currentNrOfRetries >= maxRetries {
			// skipping retry... returning the error
			info().
				Int("retry_nr", currentNrOfRetries).
				Int("max_retry", maxRetries).
				Msg(color.LightRed.Render("skipping retry..."))
			return false
		} else {
			info().
				Int("retry_nr", currentNrOfRetries).
				Int("max_retry", maxRetries).
				Msg(color.LightYellow.Render("retrying..."))
			currentNrOfRetries++

			select {
			case <-r.ctx.Done():
			case <-time.After(time.Millisecond * time.Duration(retryDelayMs)):
			}

			//time.Sleep(time.Millisecond * time.Duration(retryDelayMs))
			return true
		}
	}

retryCheck:
	// Check if it's not terminating...
	if r.IsTerminating() {
		return nil, codes.ErrTerminating
	}

	// Resolve...
	if op == Read && r.nrOfActiveReplicas.Get() >= 1 {
		info().Msg(color.LightGreen.Render("using `replicas` as active connection pool"))
		r.activeReplicasLock.RLock()
		_resolveOptions.connPool = r.activeReplicas
		_resolveOptions.cachedNrOfConnections = r.nrOfActiveReplicas.Get()
		r.activeReplicasLock.RUnlock()
		connPool = r.policy.Resolve(_resolveOptions)
	} else if r.nrOfActiveSources.Get() >= 1 { // If there are at least 1 source, then try resolving through them!
		// There are retries and other stuff that the resolver should do before returning...
		info().Msg(color.LightGreen.Render("using `sources` as active connection pool"))
		r.activeSourcesLock.RLock()
		_resolveOptions.connPool = r.activeSources
		_resolveOptions.cachedNrOfConnections = r.nrOfActiveSources.Get()
		r.activeSourcesLock.RUnlock()
		connPool = r.policy.Resolve(_resolveOptions)
	} else { // If there is nothing... then we should give an error?!
		if retryConnection(1) {
			goto retryCheck
		}
		//  return error!
		// this error should be handled by the upper function
		// the parent function should find another resolver
		return nil, codes.ErrFailedToFindAnActiveConnectionPool
		//panic(msg)
	}

	// well, this case is when the selected connection was active (monitoring hasn't  catched it as offline...)
	// we can change the host if some of them it's available...
	// so, it's the same thing as at the upper case
	if connPool == nil {
		if retryConnection(2) {
			goto retryCheck
		}
		// should we show the statement? from security perspective it should not be shown...
		return nil, codes.ErrFailedToFindAnActiveConnectionPool
		//panic(msg)
	}

	// Now Prepare statement through the selected/found connection pool
	// TODO: what is this for?!
	if stmt.DB.PrepareStmt {
		if preparedStmt, ok := r.dbResolver.prepareStmtStore[connPool]; ok {
			info().Msg("returning the prepared statement")
			return &gorm.PreparedStmtDB{
				ConnPool: connPool,
				Mux:      preparedStmt.Mux,
				Stmts:    preparedStmt.Stmts,
			}, nil
		}
	}

	return
}

func (r *resolver) call(fc func(connPool gorm.ConnPool) error) error {
	for _, s := range r.sources {
		if err := fc(s.pool); err != nil {
			return err
		}
	}

	for _, r := range r.replicas {
		if err := fc(r.pool); err != nil {
			return err
		}
	}
	return nil
}
