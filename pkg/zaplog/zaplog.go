package zaplog

import (
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Opts struct {
	Host             string
	Service          string
	MicroService     string
	Version          string
	Debug            bool
	IsDevEnvironment bool
}

// New creates a new logger instance.
func New(writer io.Writer, opts Opts) *zap.Logger {
	var (
		conf    zapcore.EncoderConfig
		encoder zapcore.Encoder
		lvl     zapcore.Level
		core    zapcore.Core
	)

	if opts.Debug {
		lvl = zapcore.DebugLevel
	}

	if opts.IsDevEnvironment {
		conf = zap.NewDevelopmentEncoderConfig()
		conf.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(conf)
	} else {
		conf = zap.NewProductionEncoderConfig()
		encoder = zapcore.NewJSONEncoder(conf)
	}

	core = zapcore.NewCore(
		encoder,
		zapcore.AddSync(writer),
		lvl,
	)

	logger := zap.New(core, zap.AddCaller()).With(
		zap.String("version", opts.Version),
		zap.String("host", opts.Host),
		zap.String("service", opts.Service),
		zap.String("microservice", opts.MicroService))
	zap.RedirectStdLog(logger)

	return logger
}
