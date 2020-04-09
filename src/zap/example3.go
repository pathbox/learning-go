package main

import (
	"time"

	"go.uber.org/zap"
)

const url = "http://example.com"

func main() {
	logger := zap.NewExample() // or NewProduction, or NewDevelopment
	defer logger.Sync()
	// In most circumstances, use the SugaredLogger. It's 4-10x faster than most
	// other structured logging packages and has a familiar, loosely-typed API.
	sugar := logger.Sugar()
	sugar.Infow("Failed to fetch URL", "url", url, "attempt", 3, "backoff", time.Second)
	sugar.Infof("Failed to fetch URL: %s", url)

	// In the unusual situations where every microsecond matters, use the
	// Logger. It's even faster than the SugaredLogger, but only supports
	// structured logging.
	logger.Info("Failed to fetch URL.", zap.String("url", url), zap.Int("attempt", 3), zap.Duration("backoff", time.Second))
}
