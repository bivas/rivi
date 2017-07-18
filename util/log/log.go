package log

import (
	"os"

	"github.com/apex/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Meta struct {
	Key   string
	Value interface{}
}

type MetaFields []Meta

func (m *MetaFields) flat() []interface{} {
	result := make([]interface{}, 0)
	for _, item := range *m {
		result = append(result, item.Key, item.Value)
	}
	return result
}

type Logger interface {
	Debug(format string, args ...interface{})
	DebugWith(meta MetaFields, format string, args ...interface{})
	Info(format string, args ...interface{})
	InfoWith(meta MetaFields, format string, args ...interface{})
	Warning(format string, args ...interface{})
	WarningWith(meta MetaFields, format string, args ...interface{})
	Error(format string, args ...interface{})
	ErrorWith(meta MetaFields, format string, args ...interface{})

	Get(name string) Logger
}

type logger struct {
	*zap.SugaredLogger
}

func (l *logger) Debug(format string, args ...interface{}) {
	l.SugaredLogger.Debugf(format, args...)
}

func (l *logger) DebugWith(meta MetaFields, format string, args ...interface{}) {
	l.SugaredLogger.With(meta.flat()...).Debugf(format, args)
}

func (l *logger) Info(format string, args ...interface{}) {
	l.SugaredLogger.Infof(format, args...)
}

func (l *logger) InfoWith(meta MetaFields, format string, args ...interface{}) {
	l.SugaredLogger.With(meta.flat()...).Infof(format, args)
}

func (l *logger) Warning(format string, args ...interface{}) {
	l.SugaredLogger.Warnf(format, args)
}

func (l *logger) WarningWith(meta MetaFields, format string, args ...interface{}) {
	l.SugaredLogger.With(meta.flat()...).Warnf(format, args)
}

func (l *logger) Error(format string, args ...interface{}) {
	l.SugaredLogger.Errorf(format, args)
}

func (l *logger) ErrorWith(meta MetaFields, format string, args ...interface{}) {
	l.SugaredLogger.With(meta.flat()...).Errorf(format, args)
}

func (l *logger) Get(name string) Logger {
	return &logger{SugaredLogger: l.Named(name)}
}

func Build() Logger {
	newLogger := &logger{}

	level := zap.InfoLevel
	if len(os.Getenv("RIVI_DEBUG")) > 0 {
		level = zap.DebugLevel
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "name",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   func(zapcore.EntryCaller, zapcore.PrimitiveArrayEncoder) {},
	}

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       true,
		Encoding:          "console",
		EncoderConfig:     encoderConfig,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stdout"},
		DisableStacktrace: true,
	}

	zapLogger, err := config.Build()
	if err != nil {
		log.Errorf("Unable to create logger. Using NoOpLogger. %s", err)
		return &NoOpLogger{}
	}

	newLogger.SugaredLogger = zapLogger.Sugar().Named("rivi")

	return newLogger
}
