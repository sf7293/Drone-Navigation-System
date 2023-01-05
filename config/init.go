package config

import (
	"os"

	core_config "github.com/sf7293/Drone-Navigation-System/config/core"
	"path/filepath"
	"runtime"
)

var (
	// Env is an "environmental variable" which shows the type of env for app to be run
	// It's accepted values are: test, dev, prod
	Env string
	// IsRunningInContainer is a boolean which indicates whether the code is running in a container or not
	// because if the code is running in container, we must get it's config from /config/file/ path instead
	IsRunningInContainer bool
	// BaseEnvConfigPath is a base config path accompanied with an environement like config/prod
	BaseEnvConfigPath string
	// BaseRootConfigPath is a the base config path no matter what the Environment is like config/
	BaseRootConfigPath string
)

func Init() {
	// core_config.EnvDev is the default value and if the "env" variable wasn't set, the Env variable would get this default value
	Env = os.Getenv("env")
	if len(Env) == 0 {
		Env = core_config.EnvDev
	}

	_, b, _, _ := runtime.Caller(0)
	ConfigBasePath := filepath.Dir(b)

	RunningInContainer := os.Getenv("running_in_container")
	if len(RunningInContainer) > 0 {
		IsRunningInContainer = true
		ConfigBasePath = "/config"
	}

	BaseRootConfigPath = ConfigBasePath + "/file"
	BaseEnvConfigPath = ConfigBasePath + "/file/dev"
	if Env == core_config.EnvProd {
		BaseEnvConfigPath = ConfigBasePath + "/file/prod"
	} else if Env == core_config.EnvTest {
		BaseEnvConfigPath = ConfigBasePath + "/file/test"
	}

	// load configs from toml files
	core_config.Loader(BaseEnvConfigPath, BaseRootConfigPath, App, Tracer, Translation, TransItem)
}
