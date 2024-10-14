package dbresolver

import (
	"sync"
	"time"

	"github.com/kyaxcorp/go-db/codes"
	"github.com/kyaxcorp/go-helper/_struct"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type PRoundRobin struct {
	PolicyPingOptions
	PolicyCommon

	currentConnIndex int
	lock             sync.Mutex
	ifStarted        bool
}

func NewRoundRobinPolicy() (*PRoundRobin, error) {
	p := &PRoundRobin{}
	_err := _struct.SetDefaultValues(p)
	if _err != nil {
		// Throw error
		return nil, codes.ErrFailedToSetDefaultValuesToRoundRobinPolicy
	}
	return p, nil
}

// Resolve -> get the connection based on the current policy
func (p *PRoundRobin) Resolve(resolverOptions *resolveOptions) gorm.ConnPool {
	funcStartTime := time.Now().Nanosecond()
	policyName := "round_robin"
	p.setPolicyName(policyName)

	// Set the connections pool to policy object for later usage...
	// There can be a problem if they are changeable... if yes then this is not ok!
	if p.resolverOptions == nil {
		p.resolverOptions = resolverOptions
	}

	info := func() *zerolog.Event {
		return p.LInfoF(policyName + ".Resolve")
	}
	defer printConsumedTime(info, funcStartTime)

	info().Msg("resolving...")

	// TODO: we should also launch a monitorer... but it will be turned off?! We need a context...
	// TODO: maybe the monitoring should be in other place, not here in the resolver?!...

	// TODO: set defaults to getConnOptions
	// The connection resolver
	connPool, _err := getConnection(&getConnOptions{
		ConnResolver: p,
	})
	if _err != nil {
		// TODO: Failed to get connection...
	}
	return connPool
}

// GetConnID -> Get's the connection id from the stack based on the policy
func (p *PRoundRobin) GetConnID() int {
	p.lock.Lock()
	defer p.lock.Unlock()
	// if we have a single defined connection, then return the fist one

	nrOfConnections := p.resolverOptions.GetNrOfConns()
	// if there is only one connection available, we
	// will return 0 (the only one)
	if nrOfConnections == 1 {
		return 0
	} else {
		// they are multiple connections here...
		// check if a reset needs to be done

		// Check if the 0 index has being returned initially
		if !p.ifStarted {
			p.ifStarted = true
			return 0
		}

		if p.currentConnIndex+1 > nrOfConnections-1 {
			// Reset back to 0
			p.currentConnIndex = 0
		} else {
			// Increment the current connection
			p.currentConnIndex++
		}
		return p.currentConnIndex
	}
}
