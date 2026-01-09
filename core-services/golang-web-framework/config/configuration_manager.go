package config

import (
	"fmt"
	"os"
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
	CONFIG_FILE_LOCATION := os.Getenv("CONFIG_FILE_LOCATION")
	var viperInstance = viper.New()
	fmt.Printf("Start to laod configs from %s", CONFIG_FILE_LOCATION)

	viperInstance.SetConfigFile(CONFIG_FILE_LOCATION)
	viperInstance.SetConfigType("yaml")

	if err := viperInstance.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s: error details: %s",CONFIG_FILE_LOCATION,  err)
	}

	return viperInstance
}

func (manager *ConfigurationManager) GetConfigFor(configKey string) string {
	return manager.viper.GetString(configKey)
}
