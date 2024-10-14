package dbresolver

import (
	"gorm.io/gorm"
	"math/rand"
)

// It's hard to know how to spread the traffic evenly, for that we need traffic analysis by communicating
// directly with the servers... or cluster
// TODO: We should know in % how busy is a node...
// TODO: for that we should have instances that are connected with the nodes, and which will read the % very often...
type PLoadBalancing struct {
	PolicyPingOptions
	PolicyCommon
}

func (p *PLoadBalancing) Resolve(resolverOptions *resolveOptions) gorm.ConnPool {
	p.setPolicyName("load_balancing")
	return resolverOptions.connPool[rand.Intn(len(resolverOptions.connPool))].pool
}
