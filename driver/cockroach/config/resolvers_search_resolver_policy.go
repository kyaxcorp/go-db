package config

import "github.com/kyaxcorp/go-helper/conv"

type SearchForAnActiveResolverIfDownPolicy struct {
	IsEnabled string `yaml:"is_enabled" mapstructure:"is_enabled" default:"yes"`
	//
	DelayMsBetweenSearches uint16 `yaml:"delay_ms_between_searches" mapstructure:"delay_ms_between_searches" default:"2000"`
	// MaxRetries -> -1 -> infinite! Maximum nr of retries....
	MaxRetries int16 `yaml:"max_retries" mapstructure:"max_retries" default:"3"`
}

func (s *SearchForAnActiveResolverIfDownPolicy) GetIsEnabled() bool {
	return conv.ParseBool(s.IsEnabled)
}

func (s *SearchForAnActiveResolverIfDownPolicy) GetDelayMsBetweenSearches() uint16 {
	return s.DelayMsBetweenSearches
}

func (s *SearchForAnActiveResolverIfDownPolicy) GetMaxRetries() int16 {
	return s.MaxRetries
}
