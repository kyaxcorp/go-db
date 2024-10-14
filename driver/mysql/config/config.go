package config

import (
	"github.com/kyaxcorp/go-db/driver"
	"github.com/kyaxcorp/go-helper/conv"
	loggerConfig "github.com/kyaxcorp/go-logger/config"
)

type Config struct {
	IsEnabled   string `yaml:"is_enabled" mapstructure:"is_enabled" default:"yes"`
	Description string

	Credentials `yaml:"global_credentials" mapstructure:"global_credentials"`
	// These options should be on all connections!
	Charset   string `yaml:"charset" mapstructure:"charset" default:"utf8mb4"`
	ParseTime string `yaml:"parse_time" mapstructure:"parse_time" default:"yes"`

	// GORM perform write (create/update/delete) operations run inside a transaction to ensure data consistency, which
	// is bad for performance, you can disable it during initialization
	SkipDefaultTransaction string `yaml:"skip_default_transaction" mapstructure:"skip_default_transaction" default:"no"`

	// These options are only for OnConnect event!
	OnConnectOptions OnConnectOptions `yaml:"on_connect_options" mapstructure:"on_connect_options"`

	// Here are the Connection Resolvers!
	// Define as many as you need, but usually you just need 1 (if you have a cluster)
	// a resolver can be as a Region Delimiter
	Resolvers []Resolver
	// Here we define reconnection things and other values
	SearchForAnActiveResolverIfDownPolicy `yaml:"search_for_an_active_resolver_if_down_policy" mapstructure:"search_for_an_active_resolver_if_down_policy"`
	// This is the logger configuration!
	Logger loggerConfig.Config
}

func (c *Config) GetDbName() string {
	return c.DbName
}

func (c *Config) GetDbUser() string {
	return c.User
}

func (c *Config) GetDbType() string {
	return "mysql"
}

func (c *Config) GetLogger() *loggerConfig.Config {
	return &c.Logger
}

func (c *Config) GetIsEnabled() bool {
	return conv.ParseBool(c.IsEnabled)
}

func (c *Config) GetResolvers() []driver.Resolver {
	var resolvers []driver.Resolver
	for _, v := range c.Resolvers {
		resolvers = append(resolvers, &v)
	}
	return resolvers
	//return c.Resolvers
}

func (c *Config) GetOnConnectOptions() driver.ConfigOnConnectOptions {
	return &c.OnConnectOptions
}

func (c *Config) GetSelf() driver.Config {
	return c
}

func (c *Config) GetSkipDefaultTransaction() bool {
	return conv.ParseBool(c.SkipDefaultTransaction)
}

func (c *Config) GetSearchForAnActiveResolverIfDownPolicy() driver.SearchForAnActiveResolverIfDownPolicy {
	return &c.SearchForAnActiveResolverIfDownPolicy
}
