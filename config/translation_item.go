package config

import (
	config "github.com/sf7293/Drone-Navigation-System/config/core"
)

type TranslationItem struct {
	Label       string `mapstructure:"label"`
	Description string `mapstructure:"description"`
	Translation string `mapstructure:"translation"`
}

type TranslationItemConfig struct {
	config.TomlInterface
	Errors []*TranslationItem `mapstructure:"errors"`
	Labels []*TranslationItem `mapstructure:"labels"`
}

var TransItem TranslationItemConfig

// a map from [language][translation_label] to Translation struct
var TranslationMap = map[string]map[string]*TranslationItem{}

func (s TranslationItemConfig) Load() {
	var transConfigs = map[string]*TranslationItemConfig{}

	activeLangCodes := Translation.ActiveLanguages
	for _, langCode := range activeLangCodes {
		transConfig := TranslationItemConfig{}
		config.LoadConfig("translation/"+langCode, &transConfig)
		updateTranslationMap(&transConfig, langCode)
		transConfigs[langCode] = &transConfig
	}
}

func updateTranslationMap(transConfig *TranslationItemConfig, langCode string) {
	translations := []*TranslationItem{}
	translations = append(translations, transConfig.Errors...)
	translations = append(translations, transConfig.Labels...)

	for _, trans := range translations {
		//logger.ZSLogger.Infow("translation item loaded", "lang_code", langCode, "translation_label", trans.Label)
		//_ = logger.ZSLogger.Sync()

		mapLabels, ok := TranslationMap[langCode]
		if !ok {
			mapLabels = map[string]*TranslationItem{}
			TranslationMap[langCode] = mapLabels
		}

		mapLabels[trans.Label] = trans
		TranslationMap[langCode] = mapLabels
	}
}
