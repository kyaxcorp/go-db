package config

type Credentials struct {
	// Username,Password, DbName are the same for the entire cluster, but it depends if you are using some kind of
	// temporary password solution which changes every specific time... that's when you will need set the credentials
	// per host/node/server
	Host     string `yaml:"host" mapstructure:"host" default:"localhost"`
	Port     int    `yaml:"port" mapstructure:"port" default:"3306"`
	User     string `yaml:"user" mapstructure:"user" default:""`
	Password string `yaml:"password" mapstructure:"password" default:""`
	DbName   string `yaml:"db_name" mapstructure:"db_name" default:""`
}
