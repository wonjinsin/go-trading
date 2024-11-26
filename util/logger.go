package util

import (
	"context"
	"encoding/json"

	"github.com/juju/errors"
	"go.uber.org/zap"
)

var zLog *Logger

// Logger ...
type Logger struct {
	*zap.SugaredLogger
}

// With ...
func (z *Logger) With(ctx context.Context) *zap.SugaredLogger {
	if ctx == nil {
		return z.SugaredLogger
	}
	trid, ok := ctx.Value(TRID).(string)
	if !ok || trid == "" {
		return z.SugaredLogger
	}

	return z.SugaredLogger.Named(trid)
}

// NewLogger ...
func NewLogger() (*Logger, error) {
	rawJSON := []byte(`{
		"level": "debug",
		"encoding": "json",
		"outputPaths": ["stdout"],
	  	"errorOutputPaths": ["stderr"],
		"encoderConfig": {
		  "levelKey": "level",
		  "levelEncoder": "capital",
		  "timeKey": "time",
		  "timeEncoder": "iso8601",
		  "nameKey": "trid",
		  "messageKey": "msg",
		  "callerKey": "caller",
		  "functionKey": "func",
		  "callerEncoder": "short"
		}
	  }`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		return nil, errors.Annotate(err, "logger build failed")
	}

	logger, err := cfg.Build()
	if err != nil {
		return nil, errors.Annotate(err, "logger build failed")
	}

	defer logger.Sync()
	zLog = &Logger{logger.Sugar()}

	return zLog, nil
}
