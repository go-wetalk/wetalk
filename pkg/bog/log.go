package bog

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var l *zap.Logger

// ProvideSingleton provides singleton Logger instance.
func ProvideSingleton() *zap.Logger {
	if l == nil {
		c := zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.AddSync(os.Stdout),
			zap.DebugLevel)

		l = zap.New(c, zap.AddStacktrace(zap.ErrorLevel))
	}

	return l
}
