package config

import (
	"github.com/kyaxcorp/go-helper/conv"
)

type ReconnectOptions struct {
	IsEnabled string `yaml:"is_enabled" mapstructure:"is_enabled" default:"yes"`
	// ReconnectAfterSeconds -> 5 Seconds, after disconnect happened, the client will reconnect after
	// indicated time!
	ReconnectAfterSeconds uint16 `yaml:"reconnect_after_seconds" mapstructure:"reconnect_after_seconds" default:"5"`
	// MaxRetries -> -1 -> infinite! Maximum nr of retries....
	MaxRetries int16 `yaml:"max_retries" mapstructure:"max_retries" default:"3"`
}

func (r *ReconnectOptions) GetIsEnabled() bool {
	return conv.ParseBool(r.IsEnabled)
}

func (r *ReconnectOptions) GetReconnectAfterSeconds() uint16 {
	return r.ReconnectAfterSeconds
}

func (r *ReconnectOptions) GetMaxRetries() int16 {
	return r.MaxRetries
}
