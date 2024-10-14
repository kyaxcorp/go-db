package constructor

import (
	"context"
	"time"

	"github.com/gookit/color"
	"github.com/kyaxcorp/go-db/codes"
	"github.com/kyaxcorp/go-db/dbresolver"
	"github.com/kyaxcorp/go-db/driver"
	dbLogger "github.com/kyaxcorp/go-db/logger"
	"github.com/kyaxcorp/go-helper/_context"
	"github.com/kyaxcorp/go-helper/_struct"
	"github.com/kyaxcorp/go-helper/conv"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func NewClient(
	ctx context.Context,
	c driver.Config,
) (*gorm.DB, error) {
	// Creating the custom logger for the database client
	// the logger will be from the system!

	l := dbLogger.NewLogger(
		c.GetDbUser()+"-"+c.GetDbName(),
		c.GetDbType(),  // TODO: change
		*c.GetLogger(), // TODO: check if should be changed to pointer in NewLogger...!?
	)
	info := func() *zerolog.Event {
		return l.Logger.InfoF("NewClient")
	}
	warn := func() *zerolog.Event {
		return l.Logger.WarnF("NewClient")
	}
	_error := func() *zerolog.Event {
		return l.Logger.ErrorF("NewClient")
	}

	// Setting as reference for later use!
	c.GetLogger().Logger = l.Logger

	info().Msg("entering...")
	defer info().Msg("leaving...")

	// Check if it's enabled...
	if !conv.ParseBool(c.GetIsEnabled()) {
		_error().Err(codes.ErrClientIsDisabled).Msg("")
		return nil, codes.ErrClientIsDisabled
	}

	info().Msg("opening connection")

	// TODO: we should check if the main connection is ok!
	// The main connection or source can give an error if not connected
	// After we have successfully connected to the main source
	// we should use the other sources,replicas
	// Or we should simply ignore the error?

	// Loop through the available connections, in  this case we take in count
	// only the sources -> WriteRead conns

	if ctx == nil {
		ctx = _context.GetDefaultContext()
	}

	onConnectOptions := c.GetOnConnectOptions()
	//onConnectOptions := c.GetOnConnectOptions().(driver.ConfigOnConnectOptions)
	resolvers := c.GetResolvers()

	var retriedTimes int8 = 0
retryConnect:
	isConnected := false
	var db *gorm.DB
	var _err error
	isLastResolver := false
	// Get the nr of resolvers
	nrOfResolvers := len(resolvers)

	// Check if there are any resolvers...
	if nrOfResolvers == 0 {
		_error().Err(codes.ErrNoResolversHaveBeenFound).Msg("")
		return nil, codes.ErrNoResolversHaveBeenFound
	}

	// sources=connections
	totalNrOfSources := 0
	// loop through resolvers
	for resolverIndex, resolver := range resolvers {
		// Check if it's the last iterated resolver
		if resolverIndex+1 == nrOfResolvers {
			isLastResolver = true
		}

		isLastSource := false
		// Get the nr of Sources

		sources := resolver.GetSources()
		nrOfSources := len(sources)

		// Check if there are any defined sources
		if nrOfSources == 0 {
			info().
				Int("resolver_index", resolverIndex).
				Int("total_sources", nrOfSources).
				Msg("no sources have been found for current resolver...continuing...")
			continue
		}

		// Count how many sources are available...
		totalNrOfSources++

		// Loop through Sources/Connections
		// And try connecting... it will connect to the first one found and successfully!
		for sourceIndex, srcConnection := range sources {
			// Check if it's the last iterated Connection/Source
			if sourceIndex+1 == nrOfSources {
				isLastSource = true
			}
			// Set the logger to the connection
			// Set references (they are needed later)
			srcConnection.SetLogger(l.Logger)
			srcConnection.SetMasterConfig(c.GetSelf())
			// Get/Generate Connection Dialector
			connDialector := srcConnection.GetDialector()
			// Open the connection
			db, _err = gorm.Open(connDialector, &gorm.Config{
				SkipDefaultTransaction: c.GetSkipDefaultTransaction(),
				Logger:                 l,
			})

			if _err != nil {
				warn().
					Err(_err).
					Str("err_message", _err.Error()).
					Msg("failed to open/initialize db")
				continue
			}

			// Get the DB and check for errors!
			sqlDB, _err := db.WithContext(ctx).DB()
			if _err != nil {
				warn().
					Err(_err).
					Str("err_message", _err.Error()).
					Msg("invalid db")
				// return nil, codes.ErrInvalidDB
				// Continue to next connection/source
				continue
			}

			_err = sqlDB.PingContext(ctx)
			if _err != nil {
				warn().
					Err(_err).
					Str("err_message", _err.Error()).
					Msg("failed to ping host")
				continue
			}

			// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
			sqlDB.SetMaxIdleConns(resolver.GetMaxIdleConnections())
			// SetMaxOpenConns sets the maximum number of open connections to the database.
			sqlDB.SetMaxOpenConns(resolver.GetMaxOpenConnections())
			// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
			sqlDB.SetConnMaxLifetime(time.Second * time.Duration(resolver.GetConnectionMaxLifeTimeSeconds()))
			//
			sqlDB.SetConnMaxIdleTime(time.Second * time.Duration(resolver.GetConnectionMaxIdleTimeSeconds()))

			if _err != nil {
				//log.Println("mysql error", _err)

				// error, failed to connect, in this case, best way is to go
				// to next available connection
				// We need here the number of retries?
				// also enable/disable retry functionality
				// allow

				// Sleep a bit till next conn...
				if onConnectOptions.GetOnFailedDelayDurationBetweenConnections() != 0 && !isLastSource {
					select {
					case <-ctx.Done():
					case <-time.After(onConnectOptions.GetOnFailedDelayDurationBetweenConnections()):
					}
					//time.Sleep(onConnectOptions.GetOnFailedDelayDurationBetweenConnections())

				}
				continue
			} else {
				isConnected = true
				break
			}
		}
		if isConnected {
			break
		}
		if isLastResolver {
			// not used yet...
		}
	}

	if totalNrOfSources == 0 {
		_error().Err(codes.ErrNoSourcesHaveBeenFoundForAnyResolver).Msg("")
		return nil, codes.ErrNoSourcesHaveBeenFoundForAnyResolver
	}

	// If it's not connected
	if !isConnected {

		// print some params
		_error().
			Int8("retriedTimes", retriedTimes).
			Bool("retry_on_failed_connect", conv.ParseBool(onConnectOptions.GetRetryOnFailed())).
			Int8("nr_of_retries_on_failed", onConnectOptions.GetMaxNrOfRetries()).
			Err(_err).
			Msg("failed to connect")

		if conv.ParseBool(onConnectOptions.GetRetryOnFailed()) {
			if onConnectOptions.GetMaxNrOfRetries() == -1 {
				// This is infinite!
				info().Msg("retrying infinitely...")
			} else {
				if retriedTimes >= onConnectOptions.GetMaxNrOfRetries() {
					// Stop right here!
					info().Msg(color.LightRed.Render("skipping retry..."))
					goto skipRetry
				}
			}
			retryAfter := time.Second * time.Duration(onConnectOptions.GetRetryDelaySeconds())
			info().Dur("retry_after", retryAfter).Msg(color.LightYellow.Render("retrying connection..."))
			retriedTimes++
			time.Sleep(retryAfter)
			goto retryConnect
		}

	skipRetry:
		if conv.ParseBool(onConnectOptions.GetPanicOnFailed()) {
			msg := "failed completely to connect to database... panicking..."
			warn().Msg(msg)
			panic(msg)
		}
		_error().Err(codes.ErrFailedToConnectToDB).Msg("")
		return nil, codes.ErrFailedToConnectToDB
	}

	// if it's connected
	if isConnected {
		retriedTimes = 0

		// loop through the resolvers
		dbResolver := dbresolver.New()
		// Set main Config in the dbResolver
		dbResolver.SetMainConfig(c)

		for _, resolver := range resolvers {

			// Define connections stacks
			var srcConnections []gorm.Dialector
			var repConnections []gorm.Dialector
			// loop through resolver's sources
			for _, connection := range resolver.GetSources() {
				// set references
				connection.SetLogger(l.Logger)
				connection.SetMasterConfig(c.GetSelf()) // TODO
				// push to source connections
				// TODO: should they be unique, should we check before adding?! CHeck by checksum?!
				srcConnections = append(srcConnections, connection.GetDialector())
			}
			// loop through resolver's replicas
			for _, connection := range resolver.GetReplicas() {
				// set references
				connection.SetLogger(l.Logger)
				connection.SetMasterConfig(c.GetSelf()) // TODO
				// push to replicas connections
				// TODO: should they be unique, should we check before adding?!
				repConnections = append(repConnections, connection.GetDialector())
			}

			// log.Println("POLICY NAME", resolver.GetPolicyName())

			policyOptions := resolver.GetPolicyOptions()

			resolverReconnectOptions := resolver.GetReconnectOptions()

			var isPolicyGenerated = false
			var _policy dbresolver.Policy

			switch resolver.GetPolicyName() {
			case dbresolver.TConsecutive:
				if policyOptions == nil {
					generatedPolicy := &dbresolver.PConsecutive{}
					_err = _struct.SetDefaultValues(generatedPolicy)
					if _err != nil {
						// Throw error
						return nil, codes.ErrFailedToSetDefaultValuesToConsecutivePolicy
					}
					_policy = generatedPolicy
				}
			case dbresolver.TLoadBalancing:
				if policyOptions == nil {
					generatedPolicy := &dbresolver.PLoadBalancing{}
					_err = _struct.SetDefaultValues(generatedPolicy)
					if _err != nil {
						// Throw error
						return nil, codes.ErrFailedToSetDefaultValuesToLoadBalancingPolicy
					}
					_policy = generatedPolicy
				}
			case dbresolver.TRandom:
				if policyOptions == nil {
					generatedPolicy := &dbresolver.PRandom{}
					_err = _struct.SetDefaultValues(generatedPolicy)
					if _err != nil { // TODO:
						return nil, codes.ErrFailedToSetDefaultValuesToRandomPolicy
					}
					isPolicyGenerated = true
				}
			case dbresolver.TRoundRobin:
				if policyOptions == nil {
					_policy, _err = dbresolver.NewRoundRobinPolicy()
					if _err != nil {
						return nil, _err
					}
					isPolicyGenerated = true
				}
			default:
				return nil, codes.ErrInvalidPolicyName
			}

			if !isPolicyGenerated {
				info().Msg("policy options have been defined by the resolver")
				if policyOptions != nil {
					_policy = policyOptions.(dbresolver.Policy)
				}
			} else {
				info().Msg("policy options have been generated using default settings")
			}

			// Set policy reconnect options
			_policy.SetIsPingRetryEnabled(resolverReconnectOptions.GetIsEnabled())
			_policy.SetPingRetryTimes(resolverReconnectOptions.GetMaxRetries())
			_policy.SetPingRetryDelaySeconds(resolverReconnectOptions.GetReconnectAfterSeconds())

			// Define tables stack
			var tables []interface{}

			// Loop through resolver's tables
			for _, tableName := range resolver.GetTables() {
				// push the tables to the stack
				// TODO: should they be unique, should we check before adding?!
				tables = append(tables, tableName)
			}

			resolverConfig := dbresolver.Config{
				Sources:  srcConnections,
				Replicas: repConnections,
				Policy:   _policy,
				// Added logger here...
				Logger: l.Logger,
				Ctx:    ctx,
			}

			//info().Interface("EEEEEE", resolverConfig).Msg("e")

			// Register the resolver

			_, _err = dbResolver.Register(resolverConfig, tables...)

			if _err != nil {
				// error...
				// Well the thing is that if we received here an error about connecting
				// to a host, it means that the primary connection also failed... after we have
				// connected to it... so if we try and start from the beginning, then we could see the error
				// that is arising... also, there may be a problem with the nr of connections (opened ports)
				// from the client side...
				// so retrying from here will be enough, and the process will exit based on the
				// OnConnectOptions that are being set in the config...

				retryAfter := time.Second * time.Duration(onConnectOptions.GetRetryDelaySeconds())
				info().Dur("retry_after", retryAfter).Msg(color.LightYellow.Render("dbResolver, retrying connection..."))
				retriedTimes++
				time.Sleep(retryAfter)
				goto retryConnect
			}

			// Set some defaults to the resolver...
			dbResolver.SetConnMaxIdleTime(time.Second * time.Duration(resolver.GetConnectionMaxIdleTimeSeconds())).
				SetConnMaxLifetime(time.Second * time.Duration(resolver.GetConnectionMaxLifeTimeSeconds())).
				SetMaxIdleConns(resolver.GetMaxIdleConnections()).
				SetMaxOpenConns(resolver.GetMaxOpenConnections())

			// I don't know if we are setting context to all sources,
			// but we are doing it just to be sure...
			/*sqlDB, _err := dbResolver.WithContext(ctx).DB()
			if _err != nil {
				warn().Str("err_message", _err.Error()).Msg("invalid db")
				return nil, codes.ErrInvalidDB
			}
			// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
			sqlDB.SetMaxIdleConns(resolver.MaxIdleConnections)
			// SetMaxOpenConns sets the maximum number of open connections to the database.
			sqlDB.SetMaxOpenConns(resolver.MaxOpenConnections)
			// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
			sqlDB.SetConnMaxLifetime(time.Second * time.Duration(resolver.ConnectionMaxLifeTimeSeconds))
			//
			sqlDB.SetConnMaxIdleTime(time.Second * time.Duration(resolver.ConnectionMaxIdleTimeSeconds))*/
		}

		// Set this plugin...
		_err = db.Use(dbResolver)
		if _err != nil {
			warn().Str("err_message", _err.Error()).Msg("failed to use db resolver plugin")
			return nil, codes.ErrFailedToUseDBResolverPlugin
		}
	}

	return db, nil
}
