package redis

import (
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/extra/redisotel/v9"

	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/config"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
)

// ClientConfig is populated by the consuming service's own configuration
// loading. This package does not read connection config keys itself —
// reading config keys is the composition root's job, matching existing
// framework conventions.
type ClientConfig struct {
	Address      string
	Password     string
	DB           int
	PoolSize     int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// NewClient builds a single-node *redis.Client (no cluster/sentinel support)
// from cfg. When otel.enabled=true it instruments the client with OTel
// tracing/metrics hooks, mirroring httpclient.NewHTTPClient and
// awsclient.LoadDefaultConfig. When otel.enabled=false the client has no
// instrumentation overhead.
func NewClient(cfg ClientConfig) *goredis.Client {
	client := goredis.NewClient(&goredis.Options{
		Addr:         cfg.Address,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	if config.GetConfigurationManagerInstance().GetConfigBoolFor("otel.enabled") {
		logger := logging.GetLoggerInstanceForComponentByType(&ClientConfig{})
		if err := redisotel.InstrumentTracing(client); err != nil {
			logger.LogErrorfFor("Error instrumenting redis client tracing: %s", err)
		}
		if err := redisotel.InstrumentMetrics(client); err != nil {
			logger.LogErrorfFor("Error instrumenting redis client metrics: %s", err)
		}
	}

	return client
}
