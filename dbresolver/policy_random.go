package dbresolver

import (
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"log"
	"math/rand"
)

type PRandom struct {
	PolicyPingOptions
	PolicyCommon

	resolverOptions *resolveOptions
}

/*
 * The random resolver will work as:
 * - it will generate the random id, based on the received id, it will check if the node is responding through Ping
 *
 */

func (p *PRandom) Resolve(resolverOptions *resolveOptions) gorm.ConnPool {
	p.setPolicyName("random")

	if p.resolverOptions == nil {
		p.resolverOptions = resolverOptions
	}
	info := func() *zerolog.Event {
		return p.LInfoF("Resolve")
	}

	connPoolLen := len(resolverOptions.connPool)
	log.Println(connPoolLen)
	connPoolId := rand.Intn(connPoolLen)
	info().Int("conn_pool_len", connPoolLen).
		Int("conn_pool_id", connPoolId).
		Msg("resolving...")
	return resolverOptions.connPool[connPoolId].pool
}
