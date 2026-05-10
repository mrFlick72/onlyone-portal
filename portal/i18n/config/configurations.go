package config

import (
	"context"

	aws_config "github.com/aws/aws-sdk-go-v2/config"
	aws_s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/awsclient"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/config"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
	bundles3 "github.com/mrflick72/onlyone-portal/portal/i18n/adapter/i18n/s3"
	"github.com/mrflick72/onlyone-portal/portal/i18n/domain/i18n"
)

const (
	defaultLanguageFallback = i18n.Locale("en_en")
	defaultLanguageClaim    = "preferred_language"
	awsRegion               = "eu-central-1"
)

var configurationManager = config.GetConfigurationManagerInstance()
var logger = logging.GetLoggerInstance()

func NewBundleRepository() i18n.BundleRepository {
	cfg, err := awsclient.LoadDefaultConfig(
		context.TODO(),
		aws_config.WithRegion(awsRegion),
	)
	if err != nil {
		logger.LogErrorfFor("unable to load SDK config: %s", err.Error())
		panic("unable to load SDK config, " + err.Error())
	}
	return bundles3.NewS3BundleRepository(
		configurationManager.GetConfigFor("i18n.s3.bundle.bucket-name"),
		aws_s3.NewFromConfig(cfg),
	)
}

func NewLanguageResolver() *i18n.LanguageResolver {
	claim := configurationManager.GetConfigFor("i18n.jwt.language-claim")
	if claim == "" {
		claim = defaultLanguageClaim
	}
	defaultLang := i18n.Locale(configurationManager.GetConfigFor("i18n.default-language"))
	if defaultLang.IsEmpty() {
		defaultLang = defaultLanguageFallback
	}
	return i18n.NewLanguageResolver(claim, defaultLang)
}
