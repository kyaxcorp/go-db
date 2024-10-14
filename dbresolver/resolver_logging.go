package dbresolver

import "github.com/rs/zerolog"

// Logging
func (r *resolver) LInfoF(functionName string) *zerolog.Event {
	return r.dbResolver.LInfoF(functionName)
}

func (r *resolver) LDebugF(functionName string) *zerolog.Event {
	return r.dbResolver.LDebugF(functionName)
}

func (r *resolver) LWarnF(functionName string) *zerolog.Event {
	return r.dbResolver.LWarnF(functionName)
}

func (r *resolver) LErrorF(functionName string) *zerolog.Event {
	return r.dbResolver.LErrorF(functionName)
}
