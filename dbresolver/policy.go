package dbresolver

import (
	"gorm.io/gorm"
)

type Policy interface {
	Resolve(resolveOptions *resolveOptions) gorm.ConnPool
	// Ping retry...
	SetIsPingRetryEnabled(isEnabled bool)
	SetPingRetryTimes(retryTimes int16)
	SetPingRetryDelaySeconds(delaySeconds uint16)
}
