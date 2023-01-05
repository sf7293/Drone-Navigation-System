package config

import (
	config "github.com/sf7293/Drone-Navigation-System/config/core"
)

type TracerConfig struct {
	config.TomlInterface
	IsActive bool `mapstructure:"is_active"`
	//The address for communicating with agent via UDP
	JaegerAgentAddress string `mapstructure:"agent_address"`
	// The reporter's flush interval in seconds
	JaegerReporterFlushInterval int `mapstructure:"reporter_flush_interval"`
}

var Tracer TracerConfig

func (s TracerConfig) Load() {
	config.LoadConfig("tracer", &Tracer)
}
