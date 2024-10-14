package config

import (
	"time"

	"github.com/kyaxcorp/go-helper/conv"
)

type OnConnectOptions struct {
	RetryOnFailed string `yaml:"retry_on_failed" mapstructure:"retry_on_failed" default:"yes"`
	// -1 means -> infinite! by default it should not be infinite! The Client who send the request should handle
	// the query to retry if something!... This Retry only helps/saves from simple cases like internet connection failure
	// The DB Client should not retry till infinite, because it's not correct to save a query to memory or somewhere else
	// This is why, the logic part or the flow part should be on client side who makes the request
	MaxNrOfRetries int8 `yaml:"retry_times_on_failed" mapstructure:"retry_times_on_failed" default:"3"`
	// This is the delay between Total retries
	RetryDelaySeconds int8 `yaml:"retry_delay_seconds" mapstructure:"retry_delay_seconds" default:"2"`
	// Delay between connections, usually it's 0, because we don't need any delay between different connections
	// we need them to connect as fast as possible!
	OnFailedDelayDurationBetweenConnections time.Duration `yaml:"on_failed_delay_duration_between_connections" mapstructure:"on_failed_delay_duration_between_connections" default:"0"`

	PanicOnFailed string `yaml:"panic_on_failed" mapstructure:"panic_on_failed" default:"no"`
}

func (o *OnConnectOptions) GetOnFailedDelayDurationBetweenConnections() time.Duration {
	return o.OnFailedDelayDurationBetweenConnections
}

func (o *OnConnectOptions) GetRetryOnFailed() bool {
	return conv.ParseBool(o.RetryOnFailed)
}

func (o *OnConnectOptions) GetMaxNrOfRetries() int8 {
	return o.MaxNrOfRetries
}

func (o *OnConnectOptions) GetRetryDelaySeconds() int8 {
	return o.RetryDelaySeconds
}

func (o *OnConnectOptions) GetPanicOnFailed() bool {
	return conv.ParseBool(o.PanicOnFailed)
}
