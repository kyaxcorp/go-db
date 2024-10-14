package dbresolver

import (
	"strconv"

	"github.com/kyaxcorp/go-helper/conv"
)

type PolicyPingOptions struct {
	// Retry Times? -1 infinite... but it's not necessary!
	// Retry Delay?
	//

	IsPingRetryEnabled    string `default:"yes"`
	PingRetryTimes        int16  `default:"3"`
	PingRetryDelaySeconds uint16 `default:"5"`
}

func (p *PolicyPingOptions) SetIsPingRetryEnabled(isEnabled bool) {
	p.IsPingRetryEnabled = strconv.FormatBool(isEnabled)
}

func (p *PolicyPingOptions) SetPingRetryTimes(retryTimes int16) {
	p.PingRetryTimes = retryTimes
}

func (p *PolicyPingOptions) SetPingRetryDelaySeconds(delaySeconds uint16) {
	p.PingRetryDelaySeconds = delaySeconds
}

// GetPingRetryTimes -> how many times it should retry the ping...
func (p *PolicyPingOptions) GetPingRetryTimes() int16 {
	return p.PingRetryTimes
}

func (p *PolicyPingOptions) GetIsPingRetryEnabled() bool {
	return conv.ParseBool(p.IsPingRetryEnabled)
}

// GetPingRetryDelaySeconds -> when checking if connection is alive, this is the delay between rechecks
func (p *PolicyPingOptions) GetPingRetryDelaySeconds() uint16 {
	return p.PingRetryDelaySeconds
}
