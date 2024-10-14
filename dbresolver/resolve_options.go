package dbresolver

import (
	"context"
	"github.com/rs/zerolog"
)

type resolveOptions struct {
	context context.Context

	//connPool   []gorm.ConnPool
	connPool   []detailedConnPool
	dbResolver *DBResolver
	resolver   *resolver

	cachedNrOfConnections int
}

// GetNrOfConns -> get nr of connections (defined ones) from the pool
func (r *resolveOptions) GetNrOfConns() int {
	return r.cachedNrOfConnections
}

// GetConnPools -> get all connections
func (r *resolveOptions) GetConnPools() []detailedConnPool {
	return r.connPool
}

func (r *resolveOptions) GetContext() context.Context {
	return r.context
}

func (r *resolveOptions) GetDbResolver() *DBResolver {
	return r.dbResolver
}

// Logging
func (r *resolveOptions) LInfoF(functionName string) *zerolog.Event {
	return r.dbResolver.LInfoF(functionName)
}

func (r *resolveOptions) LDebugF(functionName string) *zerolog.Event {
	return r.dbResolver.LDebugF(functionName)
}

func (r *resolveOptions) LWarnF(functionName string) *zerolog.Event {
	return r.dbResolver.LWarnF(functionName)
}

func (r *resolveOptions) LErrorF(functionName string) *zerolog.Event {
	return r.dbResolver.LErrorF(functionName)
}
