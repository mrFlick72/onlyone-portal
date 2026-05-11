package config

import (
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/config"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
)

const (
	defaultLanguageFallback = i18n - api.Locale("en_en")
	defaultLanguageClaim    = "preferred_language"
	awsRegion               = "eu-central-1"
)

var configurationManager = config.GetConfigurationManagerInstance()
var logger = logging.GetLoggerInstance()

func NewBundleRepository() i18n-api.BundleRepository {
cfg, err := awsclient.LoadDefaultConfig(
context.TODO(),
aws_config.WithRegion(awsRegion),
)
if err != nil {
logger.LogErrorfFor("unable to load SDK config: %s", err.Error())
panic("unable to load SDK config, " + err.Error())
}
return bundles3.NewS3BundleRepository(
configurationManager.GetConfigFor("i18n-api.s3.bundle.bucket-name"),
aws_s3.NewFromConfig(cfg),
)
}

func NewLanguageResolver() *i18n-api.LanguageResolver {
claim := configurationManager.GetConfigFor("i18n-api.jwt.language-claim")
if claim == "" {
claim = defaultLanguageClaim
}
defaultLang := i18n-api.Locale(configurationManager.GetConfigFor("i18n-api.default-language"))
if defaultLang.IsEmpty() {
defaultLang = defaultLanguageFallback
}
return i18n-api.NewLanguageResolver(claim, defaultLang)
}
