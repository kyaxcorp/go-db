package dbresolver

import (
	"time"

	"github.com/gookit/color"
	"github.com/kyaxcorp/go-db/codes"
	"github.com/kyaxcorp/go-helper/conv"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func (dr *DBResolver) getAnActiveMaster() (*resolver, error) {
	funcStartTime := time.Now().Nanosecond()
	info := func() *zerolog.Event {
		return dr.LInfoF("dbr.getAnActiveMaster")
	}
	defer printConsumedTime(info, funcStartTime)

	// Wait until the monitoring is active
	if dr.isMonitoringActive.WaitUntilTrue().Get() == false {
		// maybe it's been canceled!
		if dr.IsTerminating() {
			return nil, codes.ErrTerminating
		}
	}

	// Lock before checking...
	dr.activeMastersLock.RLock()
	defer dr.activeMastersLock.RUnlock()
	if len(dr.activeMasters) > 0 {
		return dr.activeMasters[0], nil
	}

	return nil, codes.ErrNoActiveMasters
}

func (dr *DBResolver) resolve(stmt *gorm.Statement, op Operation) gorm.ConnPool {
	resolveStartTime := time.Now().Nanosecond()

	// TODO: we should also check the resolver for availability!?
	// Should it be checked here?!

	// TODO: how right the logic is here?
	// Before switching we should check if a resolver it's fully available to process the
	// query...if it's not available check the others for compatibility

	// Check if there are any resolvers...

	info := func() *zerolog.Event {
		return dr.LInfoF("dbr.resolve")
	}
	warn := func() *zerolog.Event {
		return dr.LWarnF("dbr.resolve")
	}
	_error := func() *zerolog.Event {
		return dr.LErrorF("dbr.resolve")
	}
	info().Msg("resolving...")

	defer func() {
		info().Msg("leaving resolver...")
		resolveEndTime := time.Now().Nanosecond()
		totalResolvedTime := resolveEndTime - resolveStartTime
		info().
			Int("resolved_in_nano", totalResolvedTime).
			Msg(color.LightBlue.Render("resolved in time -> Nano:" + conv.IntToStr(totalResolvedTime)))
	}()

	/*var resolversNames []string
	for name, _ := range dr.resolvers {
		resolversNames = append(resolversNames, name)
	}

	info().Interface("resolvers_names", resolversNames).
		Int("length", len(dr.resolvers)).
		Msg("resolvers names")*/

	var connPool gorm.ConnPool
	var _err error

	searchActiveResolverPolicy := dr.mainConfig.GetSearchForAnActiveResolverIfDownPolicy()
	searchActiveResolverCurrentRetries := 0

searchForAnActiveResolver:

	if dr.IsTerminating() {
		_error().Err(codes.ErrTerminating).Msg("")
		panic(codes.ErrTerminating.Error())
	}

	var caseNr = 0
	// if it passed the resolving process
	var resolveProcessed = false
	if len(dr.resolvers) > 0 {

		//----------------CUSTOM RESOLVERS-------------------\\
		// They'll search by table names or other possible defined options by the programmer's logic...
		// in future, more logic can be added...
		if u, ok := stmt.Clauses[usingName].Expression.(using); ok && u.Use != "" {
			if r, ok := dr.resolvers[u.Use]; ok &&
				(r.CanProcessWriteOp() || (r.CanProcessReadOp() && op == Read)) {
				info().
					Str("resolve_type", "expression_using").
					Str("use_name", u.Use).
					Msg(color.LightGreen.Render("resolving with 'expression using'"))
				connPool, _err = r.resolve(stmt, op)
				resolveProcessed = true
				caseNr = 1
			}
		}

		// Check if statement has a defined table name
		// We will check through the defined earlier TABLE->RESOLVER MAP
		if stmt.Table != "" {
			// check if there is a resolver with this defined table name
			if r, ok := dr.resolvers[stmt.Table]; ok &&
				(r.CanProcessWriteOp() || (r.CanProcessReadOp() && op == Read)) {
				info().
					Str("resolve_type", "statement_table_name").
					Str("table_name", stmt.Table).
					Msg(color.LightGreen.Render("resolving with 'statement table name'"))
				connPool, _err = r.resolve(stmt, op)
				resolveProcessed = true
				caseNr = 2
			}
		}

		// Check if statement has a defined schema
		// We will check through the defined earlier TABLE->RESOLVER MAP
		if stmt.Schema != nil {
			// check if there is a resolver by the defined schema table name
			if r, ok := dr.resolvers[stmt.Schema.Table]; ok &&
				(r.CanProcessWriteOp() || (r.CanProcessReadOp() && op == Read)) {
				info().
					Str("resolve_type", "schema_table_name").
					Str("table_name", stmt.Table).
					Msg(color.LightGreen.Render("resolving with 'schema table name'"))
				connPool, _err = r.resolve(stmt, op)
				resolveProcessed = true
				caseNr = 3
			}
		}

		// Check if statement has a defined RAW Query
		// We will check through the defined earlier TABLE->RESOLVER MAP
		if rawSQL := stmt.SQL.String(); rawSQL != "" {
			// Check if there is a resolver with the defined raw table name
			if r, ok := dr.resolvers[getTableFromRawSQL(rawSQL)]; ok &&
				(r.CanProcessWriteOp() || (r.CanProcessReadOp() && op == Read)) {
				info().
					Str("resolve_type", "raw_query").
					Msg(color.LightGreen.Render("resolving with 'raw query'"))
				connPool, _err = r.resolve(stmt, op)
				resolveProcessed = true
				caseNr = 4
			}
		}
		//----------------CUSTOM RESOLVERS-------------------\\

	}

	// Get an active master resolver
	activeMaster, _err := dr.getAnActiveMaster()

	// CHeck if termination is called...
	if _err == codes.ErrTerminating {
		msg := "db client termination signal has been called"
		warn().
			Msg(color.LightRed.Render(msg))
		// TODO: should we return any error/message?
		panic(msg)
	}

	// check if there is an active master found...
	if activeMaster != nil {
		info().
			Str("resolve_type", "global_resolver").
			Msg(color.LightGreen.Render("resolving with 'global resolver'"))
		connPool, _err = activeMaster.resolve(stmt, op)
		resolveProcessed = true
		caseNr = 5
	}

	// if it's an error...
	if _err != nil {
		// TODO: sa vad ce eroare e...
		switch _err {
		case codes.ErrFailedToFindAnActiveConnectionPool, codes.ErrNoActiveMasters:
			// in this case we should try and seek another active
			_error().
				Int("case_nr", caseNr).
				Err(_err).
				Msg(color.LightRed.Render("failed to resolve, trying to find an active resolver"))

			// TODO: in this case we should try and seek another active

			if searchActiveResolverPolicy.GetIsEnabled() {
				searchMaxRetries := searchActiveResolverPolicy.GetMaxRetries()

				if searchActiveResolverCurrentRetries >= int(searchMaxRetries) && searchMaxRetries != -1 {
					_error().Err(codes.ErrRetryTimesExhaustedForSearchActiveResolvers).Msg("")
					// We throw here a panic because there are no solutions at the moment...
					// TODO the only way to get out is to throw it as a panic...?
					panic(codes.ErrRetryTimesExhaustedForSearchActiveResolvers.Error())
				} else {
					delayMs := searchActiveResolverPolicy.GetDelayMsBetweenSearches()

					time.Sleep(time.Millisecond * time.Duration(delayMs))

					searchActiveResolverCurrentRetries++
					goto searchForAnActiveResolver
				}
			}

		}
	} else {
		// no error...

		if resolveProcessed {
			//...
		}

		// if there is no connection pool found
		// it can get here only if we don't capture a specific error type
		if connPool == nil {
			_error().Err(codes.ErrConnPoolNilNoConnectionFound).Msg("")
			panic(codes.ErrConnPoolNilNoConnectionFound.Error())
		}
		// let's search for an active resolver... (loop through all of them)

		info().
			Int("case_nr", caseNr).
			Msg(color.LightGreen.Render("returning connection pool"))
		return connPool
	}

	// return the default statement connection pool
	// (which is the first connection that the client initialized)
	info().
		Str("resolve_type", "statement_conn_pool").
		Msg(color.LightGreen.Render("returning default statement connection pool"))
	return stmt.ConnPool
}
