package dbresolver

import (
	"context"
	"time"

	"github.com/gookit/color"
	"github.com/kyaxcorp/go-helper/errors2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

/*
TODO: for not wasting time with retries or things like that, it's important to have goroutines that will always
do checkups on the databases! That means that they will handle the pings every 1 second or something like that...
And they will mark  the ConnPool available or unavailable...
Doesn't matter what policies are there... if a connection is not working, it should be remove from the stack temporarily

*/

type connResolver interface {
	GetConnID() int
	GetNrOfConns() int
	GetConnPools() []detailedConnPool
	GetPingRetryTimes() int16
	GetIsPingRetryEnabled() bool
	GetPingRetryDelaySeconds() uint16
	GetContext() context.Context

	// Logging functions
	LDebugF(functionName string) *zerolog.Event
	LInfoF(functionName string) *zerolog.Event
	LErrorF(functionName string) *zerolog.Event
	LWarnF(functionName string) *zerolog.Event
}

func checkConnection(connPool gorm.ConnPool, ctx context.Context) error {
	// TODO: PingContext should be used instead!
	if pinger, ok := connPool.(interface {
		PingContext(ctx context.Context) error
	}); ok {
		//log.Println("calling ping context")
		return pinger.PingContext(ctx)
	}
	// It doesn't have a ping mechanism...
	return nil
}

type getConnOptions struct {
	ConnResolver connResolver
}

func getConnection(connOptions *getConnOptions) (gorm.ConnPool, error) {
	funcStartTime := time.Now().Nanosecond()
	retryTime := 0
	// Get a new connection ID

	p := connOptions.ConnResolver
	/*info := func() *zerolog.Event {
		return p.LInfoF("getConnection")
	}*/
	info := func() *zerolog.Event {
		return p.LInfoF("getConnection")
	}
	debug := func() *zerolog.Event {
		return p.LDebugF("getConnection")
	}
	_error := func() *zerolog.Event {
		return p.LErrorF("getConnection")
	}

	defer printConsumedTime(info, funcStartTime)

	var _err error
retry:
	debug().Msg("getting connection id")
	// TODO: first, filter the available connection pools, after that get the connection id...
	connID := p.GetConnID()
	debug().Int("connection_id", connID).Msg("got connection id, checking connection...")

	// I have commented this because ping can cause a little latency...
	// besides this, we have the monitoring system...
	// Check the connection by calling PING with Context
	/*_err = checkConnection(
		p.GetConnPools()[connID].pool,
		p.GetContext(),
	)*/

	// TODO: some of the drivers may not have ping mechanism! that's why we should be attentive

	// If we have an error, it means the connection failed, is bad...
	if _err != nil {
		// TODO: as error or we should color it...?
		_error().Int("connection_id", connID).Msg("connection failed...checking if retry is available")

		isPingRetryEnabled := p.GetIsPingRetryEnabled()
		pingRetryTimes := p.GetPingRetryTimes()
		pingRetryDelaySeconds := p.GetPingRetryDelaySeconds()

		debug().Int("retry_time", retryTime).
			Bool("is_ping_retry_enabled", isPingRetryEnabled).
			Int16("ping_retry_times", pingRetryTimes).
			Uint16("ping_retry_delay_seconds", pingRetryDelaySeconds).
			Msg("trying retry...")

		// TODO: if the query it's part of a transaction (if transaction is started) then retry should not be available!
		// TODO: the connection should die instantly... and the caller should receive the error
		// TODO: the thing is that we cannot check if it's ok or not until we get here...?

		if !isPingRetryEnabled {
			errorMsg := "ping retry is disabled"
			_error().Msg(errorMsg)
			return nil, errors2.New(0, errorMsg)
		}

		if retryTime >= int(pingRetryTimes) && pingRetryTimes != -1 {
			// If we have exhausted or retry times...
			errorMsg := "retry times exhausted"
			_error().Msg(errorMsg)
			return nil, errors2.New(0, errorMsg)
		} else {
			debug().Msg("sleeping, and after that retrying...")
			// Retry again...
			retryTime++

			select {
			case <-connOptions.ConnResolver.GetContext().Done():
				// TODO: we should return in this case!!!
			case <-time.After(time.Duration(int(pingRetryDelaySeconds)) * time.Nanosecond):
			}
			//time.Sleep(time.Duration(int(pingRetryDelaySeconds)) * time.Nanosecond)
			debug().Msg("running retry...")
			goto retry
		}
	} else {
		debug().Int("connection_id", connID).Msg(color.LightGreen.Render("connection is ok"))
	}
	return p.GetConnPools()[connID].pool, nil
}
