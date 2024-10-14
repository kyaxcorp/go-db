package dbresolver

import "gorm.io/gorm"

// This is how we register/create the plugin for gorm usage

func New() *DBResolver {
	return &DBResolver{}
}

func Register(
	config Config,
	datas ...interface{}, // it's the tables names or structures
) (*DBResolver, error) {
	// Creating a new instance of DBResolver
	return (&DBResolver{}).Register(config, datas...)
}

func (dr *DBResolver) Register(
	config Config,
	datas ...interface{}, // it's the tables names or structures
) (*DBResolver, error) {
	if dr.prepareStmtStore == nil {
		dr.prepareStmtStore = map[gorm.ConnPool]*gorm.PreparedStmtDB{}
	}

	// creating the map for storing resolvers
	if dr.resolvers == nil {
		dr.resolvers = map[string]*resolver{}
	}

	// Set Logger to dbResolver
	if config.Logger != nil {
		dr.Logger = config.Logger
	}

	// Set the default policy in case of missing one
	if config.Policy == nil {
		config.Policy = &PRandom{}
	}

	// it's the tables names or structures
	config.datas = datas
	dr.configs = append(dr.configs, config)

	var _err error
	if dr.DB != nil {
		_err = dr.compileConfig(config)
	}

	return dr, _err
}
