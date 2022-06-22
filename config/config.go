package config

import (
	"fmt"
	"main/model/mode"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var Mode = mode.Server

func Init(cfgPath, cfgName string) {
	configName := os.Getenv("CONFIG_NAME")
	if configName != "" {
		cfgName = configName
	}

	m := os.Getenv("MODE")
	Mode = mode.NewFromString(m)

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.SetConfigName(cfgName)
	viper.SetConfigType("yml")
	viper.AddConfigPath(cfgPath)

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s ", err))
	}
}
