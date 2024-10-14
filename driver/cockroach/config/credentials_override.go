package config

type CredentialsOverrides struct {
	Host              string `yaml:"host" mapstructure:"host" default:"no"`
	Port              string `yaml:"port" mapstructure:"port" default:"no"`
	User              string `yaml:"user" mapstructure:"user" default:"no"`
	Password          string `yaml:"password" mapstructure:"password" default:"no"`
	DbName            string `yaml:"db_name" mapstructure:"db_name" default:"no"`
	Schema            string `yaml:"schema" mapstructure:"schema" default:"no"`
	TimeZone          string `yaml:"time_zone" mapstructure:"time_zone" default:"no"`
	SSLMode           string `yaml:"ssl_mode" mapstructure:"ssl_mode" default:"no"`
	SSLFactory        string `yaml:"ssl_factory" mapstructure:"ssl_factory" default:"no"`
	CACertificate     string `yaml:"ca_certificate" mapstructure:"ca_certificate" default:"no"`
	ClientCertificate string `yaml:"client_certificate" mapstructure:"client_certificate" default:"no"`
	ClientPrivateKey  string `yaml:"client_private_key" mapstructure:"client_private_key" default:"no"`
}
