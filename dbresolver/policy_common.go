package dbresolver

import (
	"context"
	"github.com/rs/zerolog"
)

type PolicyCommon struct {
	policyName string
	// resolveOptions is set on Resolve it's being called
	resolverOptions *resolveOptions
}

func (p *PolicyCommon) GetNrOfConns() int {
	return p.resolverOptions.GetNrOfConns()
}

// GetConnPools -> get all connections
func (p *PolicyCommon) GetConnPools() []detailedConnPool {
	return p.resolverOptions.GetConnPools()
}

func (p *PolicyCommon) GetContext() context.Context {
	return p.resolverOptions.GetContext()
}

func (p *PolicyCommon) setPolicyName(policyName string) {
	if p.policyName == "" {
		p.policyName = policyName
	}
}

// Logging
//

func (p *PolicyCommon) LInfoF(functionName string) *zerolog.Event {
	return p.resolverOptions.LInfoF(functionName).Str("policy_name", p.policyName)
}

func (p *PolicyCommon) LDebugF(functionName string) *zerolog.Event {
	return p.resolverOptions.LDebugF(functionName).Str("policy_name", p.policyName)
}

func (p *PolicyCommon) LWarnF(functionName string) *zerolog.Event {
	return p.resolverOptions.LWarnF(functionName).Str("policy_name", p.policyName)
}

func (p *PolicyCommon) LErrorF(functionName string) *zerolog.Event {
	return p.resolverOptions.LErrorF(functionName).Str("policy_name", p.policyName)
}
