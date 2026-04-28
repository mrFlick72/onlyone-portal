package config

import (
	"fmt"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/adapter/plan/db"

	frameworkConfig "github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/config"
)

func NewPlanRepository(dsn string) *db.PlanPostgresRepository {
	return &db.PlanPostgresRepository{ConnectionString: dsn}
}

func NewPostgresDSN() string {
	cm := frameworkConfig.GetConfigurationManagerInstance()
	return fmt.Sprintf(
		"host=%s dbname=%s user=%s password=%s sslmode=disable",
		cm.GetConfigFor("database.url"),
		cm.GetConfigFor("database.name"),
		cm.GetConfigFor("database.user"),
		cm.GetConfigFor("database.password"),
	)
}
