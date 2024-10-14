package config

import "github.com/kyaxcorp/go-helper/_struct"

func SetDefaults(config *Config) (*Config, error) {
	if config == nil {
		config = &Config{}
	}

	var _err error
	_err = _struct.SetDefaultValues(config)
	if _err != nil {
		return config, _err
	}

	// Set default values for logger
	_err = _struct.SetDefaultValues(&config.Logger)
	if _err != nil {
		return config, _err
	}

	// Check if 1 resolver exists...
	// if no, then create 1

	if len(config.Resolvers) == 0 {
		defaultResolver := Resolver{}
		_err = _struct.SetDefaultValues(&defaultResolver)
		if _err != nil {
			return config, _err
		}
		config.Resolvers = []Resolver{defaultResolver}
	}

	// loop through all existing resolvers and set defaults
	for nr, r := range config.Resolvers {
		_err = _struct.SetDefaultValues(&r)
		if _err != nil {
			return config, _err
		}

		// Check if the resolver has a primary/default source
		if len(r.Sources) == 0 {
			defaultConnection := Connection{}
			_err = _struct.SetDefaultValues(&defaultConnection)
			if _err != nil {
				return config, _err
			}
			// Credentials override
			_err = _struct.SetDefaultValues(&defaultConnection.CredentialsOverrides)
			if _err != nil {
				return config, _err
			}

			// Set the default Con nection
			r.Sources = []Connection{defaultConnection}
		}

		// Loop through sources
		for connNr, conn := range r.Sources {
			_err = _struct.SetDefaultValues(&conn)
			if _err != nil {
				return config, _err
			}
			// Credentials override
			_err = _struct.SetDefaultValues(&conn.CredentialsOverrides)
			if _err != nil {
				return config, _err
			}
			// Set back the connection
			r.Sources[connNr] = conn
		}

		// Loop through replicas
		for connNr, conn := range r.Replicas {
			_err = _struct.SetDefaultValues(&conn)
			if _err != nil {
				return config, _err
			}
			// Set back the connection
			r.Replicas[connNr] = conn
		}

		// Set back the resolver
		config.Resolvers[nr] = r
	}

	return config, nil
}
