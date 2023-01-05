package core

import (
	"strings"
	"time"

	"github.com/sf7293/Drone-Navigation-System/utils/logger"
	// Vendor Packages
	"github.com/spf13/viper"
)

const (
	EnvDev  = "dev"
	EnvProd = "prod"
	EnvTest = "test"
)

var (
	// BaseEnvConfigPath is a base config path accompanied with an environement like config/prod
	BaseEnvConfigPath string
	// BaseRootConfigPath is a the base config path no matter what the Environment is like config/
	BaseRootConfigPath string
)

type TomlInterface interface {
	Load()
}

func Loader(baseEnvConfigPath string, baseRootConfigPath string, structs ...TomlInterface) {
	defer func() {
		_ = logger.ZSLogger.Sync()
	}()

	BaseEnvConfigPath = baseEnvConfigPath
	BaseRootConfigPath = baseRootConfigPath

	// check structs input
	if len(structs) == 0 {
		logger.ZSLogger.Panic("config struct is empty")
		return
	}

	// call interface method
	for _, s := range structs {
		s.Load()
	}
}

func LoadConfig(configPath string, configStruct interface{}) {
	defer func() {
		_ = logger.ZSLogger.Sync()
	}()

	viperInstance := viper.New()
	fullConfigPath := BaseEnvConfigPath + "/" + configPath + ".toml"

	// translation files don't depend on the environment, so we've placed them in the root of config directory
	if strings.HasPrefix(configPath, "translation") {
		fullConfigPath = BaseRootConfigPath + "/" + configPath + ".toml"
	}

	viperInstance.SetConfigFile(fullConfigPath)
	err := viperInstance.ReadInConfig()
	if err != nil {
		logger.ZSLogger.Panicw("failed to read from file", "config_path", fullConfigPath, "error", err)
		return
	}

	err = viperInstance.Unmarshal(configStruct)
	if err != nil {
		logger.ZSLogger.Panicw("failed to unmarshal", "config_path", fullConfigPath, "error", err)
		return
	}
	go periodicReadAndUpdate(viperInstance, configPath, configStruct)

	logger.ZSLogger.Infow("config loaded successfully", "config_path", fullConfigPath)
}

func periodicReadAndUpdate(viperInstance *viper.Viper, configName string, config interface{}) {
	for {
		time.Sleep(5 * time.Second)

		err := viperInstance.ReadInConfig()
		for err != nil {
			logger.ZSLogger.Errorw("failed to read from config map in periodic read loop", "error", err)
			_ = logger.ZSLogger.Sync()
			continue
		}

		err = viperInstance.Unmarshal(config)
		if err != nil {
			logger.ZSLogger.Errorw("error while unmarshalling viper config in periodic loop", "error", err, "config_name", configName)
			_ = logger.ZSLogger.Sync()
		}
	}
}
