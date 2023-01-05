package config

import config "github.com/sf7293/Drone-Navigation-System/config/core"

type AppConfig struct {
	config.TomlInterface
	//ConfigVersion is just a number by increase of which we can detect that the system has read our latest configs
	//Because we have a fmt log in config/core.go which prints latest value of config.App.ConfigVersion each 5 seconds
	ConfigVersion int64  `mapstructure:"config_version"`
	HTTPPort      int    `mapstructure:"http_port"`
	SectorID      int64  `mapstructure:"sector_id"`
	ServerName    string `mapstructure:"server_name"`
	Mode          string `mapstructure:"mode"`
}

var App AppConfig

func (s AppConfig) Load() {
	config.LoadConfig("app", &App)
}
