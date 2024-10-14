package driver

import (
	"time"

	loggerConfig "github.com/kyaxcorp/go-logger/config"
	"github.com/kyaxcorp/go-logger/model"
	"gorm.io/gorm"
)

// this is the main configuration which involves all the other interfaces
type Config interface {
	GetDbName() string
	GetDbUser() string
	GetDbType() string // mysql, cockroachdb, postgres etc...
	GetLogger() *loggerConfig.Config
	GetIsEnabled() bool
	GetResolvers() []Resolver
	GetOnConnectOptions() ConfigOnConnectOptions
	GetSelf() Config
	GetSkipDefaultTransaction() bool
	GetSearchForAnActiveResolverIfDownPolicy() SearchForAnActiveResolverIfDownPolicy
}

// This is when first connecting...
type ConfigOnConnectOptions interface {
	GetOnFailedDelayDurationBetweenConnections() time.Duration
	GetRetryOnFailed() bool
	GetMaxNrOfRetries() int8
	GetRetryDelaySeconds() int8
	GetPanicOnFailed() bool
}

type SearchForAnActiveResolverIfDownPolicy interface {
	GetIsEnabled() bool
	GetDelayMsBetweenSearches() uint16
	GetMaxRetries() int16
}

// this is the resolver
type Resolver interface {
	GetSources() []Connection
	GetReplicas() []Connection
	GetPolicyName() string
	GetMaxIdleConnections() int
	GetMaxOpenConnections() int
	GetConnectionMaxLifeTimeSeconds() uint32
	GetConnectionMaxIdleTimeSeconds() uint32
	//GetPolicyOptions() dbresolver.Policy
	GetPolicyOptions() interface{}
	//SetPolicyOptions(policy dbresolver.Policy)
	SetPolicyOptions(policy interface{})
	GetReconnectOptions() ReconnectOptions
	GetTables() []string
}

type Connection interface {
	SetLogger(logger *model.Logger)
	SetMasterConfig(config interface{})
	GetDialector() gorm.Dialector
	GetReconnectOptions() ReconnectOptions
}

// this is the reconnect options for the resolver & connection
type ReconnectOptions interface {
	GetIsEnabled() bool
	GetReconnectAfterSeconds() uint16
	GetMaxRetries() int16
}
