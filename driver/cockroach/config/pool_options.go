package config

type ConnectionPoolOptions struct {
	// SetMaxIdleConnections sets the maximum number of connections in the idle connection pool.
	MaxIdleConnections int `yaml:"max_idle_connections" mapstructure:"max_idle_connections" default:"10"`
	// SetMaxOpenConnections sets the maximum number of open connections to the database.
	MaxOpenConnections int `yaml:"max_open_connections" mapstructure:"max_open_connections" default:"100"`
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	ConnectionMaxLifeTimeSeconds uint32 `yaml:"connection_max_life_time_seconds" mapstructure:"connection_max_life_time_seconds" default:"86400"`
	ConnectionMaxIdleTimeSeconds uint32 `yaml:"connection_max_idle_time_seconds" mapstructure:"connection_max_idle_time_seconds" default:"3600"`
}
