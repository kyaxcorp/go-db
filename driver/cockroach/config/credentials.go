package config

type Credentials struct {
	// Username,Password, DbName are the same for the entire cluster, but it depends if you are using some kind of
	// temporary password solution which changes every specific time... that's when you will need set the credentials
	// per host/node/server
	Host     string `yaml:"host" mapstructure:"host" default:"localhost"`
	Port     int    `yaml:"port" mapstructure:"port" default:"26257"`
	User     string `yaml:"user" mapstructure:"user" default:"root"`
	Password string `yaml:"password" mapstructure:"password" default:""`
	DbName   string `yaml:"db_name" mapstructure:"db_name" default:"defaultdb"`
	Schema   string `yaml:"schema" mapstructure:"schema" default:"public"`

	TimeZone string `yaml:"time_zone" mapstructure:"time_zone" default:""`

	// if you want to disable ssl, then write: disable
	SSLMode string `yaml:"ssl_mode" mapstructure:"ssl_mode" default:"require"`
	// SSLFactory - Optional
	SSLFactory string `yaml:"ssl_factory" mapstructure:"ssl_factory" default:""` // TODO: what's this
	// CACertificate - Optional
	CACertificate string `yaml:"ca_certificate" mapstructure:"ca_certificate" default:""`
	// ClientCertificate - Optional
	ClientCertificate string `yaml:"client_certificate" mapstructure:"client_certificate" default:""`
	// ClientPrivateKey - Optional
	ClientPrivateKey string `yaml:"client_private_key" mapstructure:"client_private_key" default:""`
}
