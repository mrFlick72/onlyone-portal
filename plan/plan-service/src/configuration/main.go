package configuration

import (
	"fmt"

	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/config"
)

func NewPostgresDSN() string {
	cm := config.GetConfigurationManagerInstance()
	return fmt.Sprintf(
		"host=%s dbname=%s user=%s password=%s sslmode=disable",
		cm.GetConfigFor("database.url"),
		cm.GetConfigFor("database.name"),
		cm.GetConfigFor("database.user"),
		cm.GetConfigFor("database.password"),
	)
}
