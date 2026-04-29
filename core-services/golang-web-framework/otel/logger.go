package otel

import "github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"

var tracingLogger = logging.GetLoggerInstanceForComponentByTypeName("OtelProvider")
