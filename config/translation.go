package config

import (
	config "github.com/sf7293/Drone-Navigation-System/config/core"
)

type TranslationConfig struct {
	config.TomlInterface
	ActiveLanguages []string `mapstructure:"active_languages"`
}

var Translation TranslationConfig

func (s TranslationConfig) Load() {
	config.LoadConfig("translation/translation", &Translation)
}
