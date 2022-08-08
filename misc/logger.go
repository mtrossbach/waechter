package misc

import (
	"sync"

	"go.uber.org/zap"
)

var Log *zap.SugaredLogger
var once sync.Once

func InitializeLogging() {
	once.Do(func() {
		logger, _ := zap.NewDevelopment()
		defer logger.Sync() // flushes buffer, if any
		Log = logger.Sugar()
	})
}
