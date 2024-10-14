package codes

import (
	"github.com/kyaxcorp/go-helper/errors2/define"
	"gorm.io/gorm"
)

const ErrCodeClientIsDisabled = 100
const ErrCodeFailedToConnectToDB = 101
const ErrCodeInvalidDB = 102
const ErrCodeNoDefinedConnections = 103
const ErrCodeAlreadyConnected = 104
const ErrCodeAlreadyConnecting = 105
const ErrCodeTerminating = 106
const ErrCodeRetryFailed = 107
const ErrCodeIsDisconnecting = 108
const ErrCodeNotConnected = 109
const ErrCodeIsReconnecting = 110
const ErrCodeReconnectDisabled = 111
const ErrCodeFailedToStartConnMonitoring = 112
const ErrCodeFailedToUseDBResolverPlugin = 113

const ErrCodeFailedToSetDefaultValuesToConsecutivePolicy = 114
const ErrCodeFailedToSetDefaultValuesToLoadBalancingPolicy = 115
const ErrCodeFailedToSetDefaultValuesToRandomPolicy = 116
const ErrCodeFailedToSetDefaultValuesToRoundRobinPolicy = 117
const ErrCodeInvalidPolicyName = 118

// ------------ DB RESOLVER-----------------\\
const ErrCodeFailedToFindAnActiveConnectionPool = 200

// Resolver connections monitoring
const ErrCodeFailedToGenerateResolverMonitoringGoroutine = 201
const ErrCodeFailedToRunResolverMonitoringGoroutine = 202

// DB Client Resolvers manager
const ErrCodeFailedToGenerateResolversManagerMonitoringGoroutine = 203
const ErrCodeFailedToRunResolversManagerMonitoringGoroutine = 204

const ErrCodeNoActiveMasters = 205

const ErrCodeFailedToStartConnectionsMonitoringForResolver = 206

const ErrCodeFailedToStartResolversMonitoring = 207

const ErrCodeRetryTimesExhaustedForSearchActiveResolvers = 208
const ErrCodeConnPoolNilNoConnectionFound = 209
const ErrCodeDbInstanceIsMissing = 210
const ErrCodeDefaultDbInstanceNameIsEmpty = 211
const ErrCodeDbClientInstanceNameEmpty = 212
const ErrCodeDbInstanceConfigurationIsMissing = 214
const ErrCodeDriverNoDefaultConfigFound = 215

const ErrCodeNoResolversHaveBeenFound = 216
const ErrCodeNoSourcesHaveBeenFoundForResolver = 217
const ErrCodeNoSourcesHaveBeenFoundForAnyResolver = 218

// ------------ DB RESOLVER-----------------\\

var ErrClientIsDisabled = define.Err(ErrCodeClientIsDisabled, "database client is disabled -> check your configuration for field: 'is_enabled'")
var ErrFailedToConnectToDB = define.Err(ErrCodeFailedToConnectToDB, "failed to connect to the db")
var ErrInvalidDB = define.Err(ErrCodeInvalidDB, gorm.ErrInvalidDB.Error())
var ErrNoDefinedConnections = define.Err(ErrCodeNoDefinedConnections, "no defined connections")
var ErrAlreadyConnected = define.Err(ErrCodeAlreadyConnected, "already connected")
var ErrAlreadyConnecting = define.Err(ErrCodeAlreadyConnecting, "already connecting")
var ErrTerminating = define.Err(ErrCodeTerminating, "context cancel or done has been called, terminating...")
var ErrRetryFailed = define.Err(ErrCodeRetryFailed, "retry failed...")
var ErrIsDisconnecting = define.Err(ErrCodeIsDisconnecting, "is disconnecting...")
var ErrNotConnected = define.Err(ErrCodeNotConnected, "not connected...")
var ErrIsReconnecting = define.Err(ErrCodeIsReconnecting, "is reconnecting...")
var ErrReconnectDisabled = define.Err(ErrCodeReconnectDisabled, "reconnect disabled...")
var ErrFailedToStartConnMonitoring = define.Err(ErrCodeFailedToStartConnMonitoring, "failed to start connection monitoring...")
var ErrFailedToUseDBResolverPlugin = define.Err(ErrCodeFailedToUseDBResolverPlugin, "failed to use db resolver plugin...")

// Policies
var ErrFailedToSetDefaultValuesToConsecutivePolicy = define.Err(ErrCodeFailedToSetDefaultValuesToConsecutivePolicy, "failed to set default values to 'Consecutive' policy")
var ErrFailedToSetDefaultValuesToLoadBalancingPolicy = define.Err(ErrCodeFailedToSetDefaultValuesToLoadBalancingPolicy, "failed to set default values to 'Load Balancing' policy")
var ErrFailedToSetDefaultValuesToRandomPolicy = define.Err(ErrCodeFailedToSetDefaultValuesToRandomPolicy, "failed to set default values to 'Random' policy")
var ErrFailedToSetDefaultValuesToRoundRobinPolicy = define.Err(ErrCodeFailedToSetDefaultValuesToRoundRobinPolicy, "failed to set default values to 'Round Robin' policy")
var ErrInvalidPolicyName = define.Err(ErrCodeInvalidPolicyName, "invalid policy name")

var ErrFailedToFindAnActiveConnectionPool = define.Err(ErrCodeFailedToFindAnActiveConnectionPool, "failed to find an active/restored connection pool")

// Resolver connections monitoring
var ErrFailedToGenerateResolverMonitoringGoroutine = define.Err(ErrCodeFailedToGenerateResolverMonitoringGoroutine, "failed to generate resolver connection monitoring goroutine")
var ErrFailedToRunResolverMonitoringGoroutine = define.Err(ErrCodeFailedToRunResolverMonitoringGoroutine, "failed to run resolver connection monitoring goroutine")

// DB Client Resolvers manager
var ErrFailedToGenerateResolversManagerMonitoringGoroutine = define.Err(ErrCodeFailedToGenerateResolversManagerMonitoringGoroutine, "failed to generate resolvers manager monitoring goroutine")
var ErrFailedToRunResolversManagerMonitoringGoroutine = define.Err(ErrCodeFailedToRunResolversManagerMonitoringGoroutine, "failed to run resolvers manager monitoring goroutine")

var ErrNoActiveMasters = define.Err(ErrCodeNoActiveMasters, "no active master resolvers")

var ErrFailedToStartConnectionsMonitoringForResolver = define.Err(ErrCodeFailedToStartConnectionsMonitoringForResolver, "failed to start connections monitoring for resolver")

var ErrFailedToStartResolversMonitoring = define.Err(ErrCodeFailedToStartResolversMonitoring, "failed to start resolvers monitoring")
var ErrRetryTimesExhaustedForSearchActiveResolvers = define.Err(ErrCodeRetryTimesExhaustedForSearchActiveResolvers, "retry times exhausted for search active resolvers")
var ErrConnPoolNilNoConnectionFound = define.Err(ErrCodeConnPoolNilNoConnectionFound, "no active masters (db resolvers), check error types...")

var ErrDbInstanceIsMissing = define.Err(ErrCodeDbInstanceIsMissing, "db instance is missing")
var ErrDefaultDbInstanceNameIsEmpty = define.Err(ErrCodeDefaultDbInstanceNameIsEmpty, "default db instance name is empty")
var ErrDbClientInstanceNameEmpty = define.Err(ErrCodeDbClientInstanceNameEmpty, "db client instance name is empty")
var ErrDbInstanceConfigurationIsMissing = define.Err(ErrCodeDbInstanceConfigurationIsMissing, "db client instance configuration is missing")
var ErrDriverNoDefaultConfigFound = define.Err(ErrCodeDriverNoDefaultConfigFound, "driver hasn't any default config")

var ErrNoResolversHaveBeenFound = define.Err(ErrCodeNoResolversHaveBeenFound, "no resolvers have have been defined/found in the config")
var ErrNoSourcesHaveBeenFoundForResolver = define.Err(ErrCodeNoSourcesHaveBeenFoundForResolver, "no sources have been defined/found for the resolver in the config")
var ErrNoSourcesHaveBeenFoundForAnyResolver = define.Err(ErrCodeNoSourcesHaveBeenFoundForAnyResolver, "no sources have been defined/found for any resolver in the config")
