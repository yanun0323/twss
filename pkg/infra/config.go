package infra

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

var (
	_Once = &sync.Once{}
)

func Init(cfgName string) error {
	var err error
	_Once.Do(
		func() {
			if cfgName == "" {
				cfgName = "config"
			}
			_, f, _, _ := runtime.Caller(0)
			cfgPath := filepath.Join(filepath.Dir(f), "../../config")
			viper.AddConfigPath(cfgPath)
			viper.AddConfigPath("./config")
			configName := os.Getenv("CONFIG_NAME")
			if configName != "" {
				cfgName = configName
			}

			viper.SetConfigName(cfgName)
			viper.AutomaticEnv()
			viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
			viper.SetConfigType("yaml")

			err = viper.ReadInConfig()
			if err != nil {
				return
			}
		},
	)
	return err
}
