package telegofx

import (
	"fmt"

	"github.com/mymmrac/telego"
	"go.uber.org/zap"
)

type zapLogger struct {
	logger *zap.Logger
}

// Debugf implements telego.Logger.
func (z *zapLogger) Debugf(format string, args ...any) {
	z.logger.Debug(fmt.Sprintf(format, args...), zap.String("format", format), zap.Any("args", args))
}

// Errorf implements telego.Logger.
func (z *zapLogger) Errorf(format string, args ...any) {
	z.logger.Error(fmt.Sprintf(format, args...), zap.String("format", format), zap.Any("args", args))
}

var _ telego.Logger = (*zapLogger)(nil)
