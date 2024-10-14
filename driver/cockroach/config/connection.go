package config

import (
	"github.com/kyaxcorp/go-db/driver"
	"github.com/kyaxcorp/go-logger/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Connection struct {
	CredentialsOverrides CredentialsOverrides `yaml:"credentials_overrides" mapstructure:"credentials_overrides"`
	Credentials          `yaml:"credentials" mapstructure:"credentials"`

	// This is only for this connection!
	ReconnectOptions `yaml:"reconnect_options" mapstructure:"reconnect_options"`

	//
	logger       *model.Logger
	masterConfig *Config
}

func (c *Connection) SetLogger(logger *model.Logger) {
	c.logger = logger
}

func (c *Connection) SetMasterConfig(config interface{}) {
	c.masterConfig = config.(*Config)
}

func (c *Connection) GetDialector() gorm.Dialector {
	dsn := c.GenerateDSN()
	c.logger.Info().
		Str("type", "postgresql").
		Str("dsn", dsn.Secured).
		Msg("generating PostgreSQL DSN")
	return postgres.Open(dsn.Plain)
}

func (c *Connection) GetReconnectOptions() driver.ReconnectOptions {
	return &c.ReconnectOptions
}
