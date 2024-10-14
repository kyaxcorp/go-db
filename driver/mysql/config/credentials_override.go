package config

type CredentialsOverrides struct {
	Host     string `yaml:"host" mapstructure:"host" default:"no"`
	Port     string `yaml:"port" mapstructure:"port" default:"no"`
	User     string `yaml:"user" mapstructure:"user" default:"no"`
	Password string `yaml:"password" mapstructure:"password" default:"no"`
	DbName   string `yaml:"db_name" mapstructure:"db_name" default:"no"`
}
