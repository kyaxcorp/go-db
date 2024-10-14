package config

import (
	"github.com/kyaxcorp/go-db/dbresolver"
	"github.com/kyaxcorp/go-db/driver"
)

// Multiple resolvers should be used on your own risk...
type Resolver struct {
	// policy... can be the name of the policy... like

	// *Random -> it used rand function...
	// *Load Balancing -> One by one, meaning that it will not trigger the same db, it will go to the next available one
	// by having statistics of usage of each DB or cluster...
	// *RoundRobin -> it will trigger in a consecutive order reaching the end and starting again from beginning
	// *Consecutive -> it will start from the beginning, if it works, then it uses that connection always... if it doesn't
	// it will go to next one until it reaches an ok connection!

	// if it's not working (depending on the policy type and algo of selection)
	// it will go to next one...if none of them working it will start from beginning or it will trigger
	// an error...? we should define how many times or how much time it will check for the connections to be ok
	PolicyName string `yaml:"policy_name" mapstructure:"policy_name" default:"round_robin"`

	//PolicyOptions PolicyOptions `yaml:"policy_options" mapstructure:"policy_options"`
	// PolicyOptions should not be exported to yaml
	// It should be used only internal as code...
	PolicyOptions dbresolver.Policy `yaml:"-"`

	ReconnectOptions `yaml:"reconnect_options" mapstructure:"reconnect_options"`

	// For which tables it's referred to...it will switch automatically the connections based on this
	// Is using this, the Root Resolver will switch/use the resolver based on this
	Tables []string
	// TODO: we can also add here Regions, so based on a TABLE COLUMN,but for that we should scan the query for words

	ConnectionPoolOptions `yaml:"connection_pool_options" mapstructure:"connection_pool_options"`

	// When writing or read writing it will always take the connections from here...
	Sources []Connection
	// When only reading, it will always take the connections from here!
	Replicas []Connection
}

func (r *Resolver) GetSources() []driver.Connection {
	var connections []driver.Connection
	for _, v := range r.Sources {
		connections = append(connections, &v)
	}
	return connections
	//return r.Sources
}

func (r *Resolver) GetReplicas() []driver.Connection {
	var connections []driver.Connection
	for _, v := range r.Replicas {
		connections = append(connections, &v)
	}
	return connections
	//return r.Replicas
}

func (r *Resolver) GetPolicyName() string {
	return r.PolicyName
}

func (r *Resolver) GetMaxIdleConnections() int {
	return r.MaxIdleConnections
}
func (r *Resolver) GetMaxOpenConnections() int {
	return r.MaxOpenConnections
}

func (r *Resolver) GetConnectionMaxLifeTimeSeconds() uint32 {
	return r.ConnectionMaxLifeTimeSeconds
}

func (r *Resolver) GetConnectionMaxIdleTimeSeconds() uint32 {
	return r.ConnectionMaxIdleTimeSeconds
}

// func (r *Resolver) GetPolicyOptions() dbresolver.Policy {
func (r *Resolver) GetPolicyOptions() interface{} {
	return r.PolicyOptions
}

// func (r *Resolver) SetPolicyOptions(policy dbresolver.Policy) {
func (r *Resolver) SetPolicyOptions(policy interface{}) {
	r.PolicyOptions = policy.(dbresolver.Policy)
}

func (r *Resolver) GetReconnectOptions() driver.ReconnectOptions {
	return &r.ReconnectOptions
}

func (r *Resolver) GetTables() []string {
	return r.Tables
}
