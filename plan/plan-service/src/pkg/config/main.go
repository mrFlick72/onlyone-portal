package config

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

var (
	configurationManager *ConfigurationManager
	once                 sync.Once
)

type ConfigurationManager struct {
	viper *viper.Viper
}

func GetConfigurationManagerInstance() *ConfigurationManager {
	once.Do(func() {
		configurationManager = &ConfigurationManager{
			viper: initViperInstance(),
		}
	})

	return configurationManager
}

func initViperInstance() *viper.Viper {
	var viperInstance = viper.New()
	fmt.Printf("Start to laod configs ")

	viperInstance.SetConfigFile("/home/vvaudi/onlyoneportal-products/applications/todo/todo-service/.env")
	viperInstance.SetConfigType("env")

	if err := viperInstance.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	return viperInstance
}

func (manager *ConfigurationManager) GetConfigFor(configKey string) string {
	return manager.viper.GetString(configKey)
}
