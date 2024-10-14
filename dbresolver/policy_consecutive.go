package dbresolver

import (
	"gorm.io/gorm"
	"math/rand"
)

type PConsecutive struct {
	PolicyPingOptions
	PolicyCommon
}

/*
 * The consecutive policy will always check in a consecutive way the availability of the nodes, meaning:
 * - if the first node is available, it will query until it's off
 * - if the first not responding to ping, it will go to next one
 */

func (p *PConsecutive) Resolve(resolverOptions *resolveOptions) gorm.ConnPool {
	p.setPolicyName("consecutive")
	return resolverOptions.connPool[rand.Intn(len(resolverOptions.connPool))].pool
}
