package thriftlion

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"go.pedge.io/lion"
)

var (
	discardLevelLoggerInstance = newDiscardLevelLogger()
)

type logger struct {
	lion.Logger
	l lion.Level
}

func newLogger(delegate lion.Logger) *logger {
	return &logger{delegate, delegate.Level()}
}

func (l *logger) AtLevel(level lion.Level) Logger {
	return newLogger(l.Logger.AtLevel(level))
}

func (l *logger) WithField(key string, value interface{}) Logger {
	return newLogger(l.Logger.WithField(key, value))
}

func (l *logger) WithFields(fields map[string]interface{}) Logger {
	return newLogger(l.Logger.WithFields(fields))
}

func (l *logger) WithKeyValues(keyValues ...interface{}) Logger {
	return newLogger(l.Logger.WithKeyValues(keyValues...))
}

func (l *logger) WithContext(context thrift.TStruct) Logger {
	return newLogger(l.WithEntryMessageContext(newEntryMessage(context)))
}

func (l *logger) Debug(event thrift.TStruct) {
	if lion.LevelDebug < l.l {
		return
	}
	l.LogEntryMessage(lion.LevelDebug, newEntryMessage(event))
}

func (l *logger) Info(event thrift.TStruct) {
	if lion.LevelInfo < l.l {
		return
	}
	l.LogEntryMessage(lion.LevelInfo, newEntryMessage(event))
}

func (l *logger) Warn(event thrift.TStruct) {
	if lion.LevelWarn < l.l {
		return
	}
	l.LogEntryMessage(lion.LevelWarn, newEntryMessage(event))
}

func (l *logger) Error(event thrift.TStruct) {
	if lion.LevelError < l.l {
		return
	}
	l.LogEntryMessage(lion.LevelError, newEntryMessage(event))
}

func (l *logger) Fatal(event thrift.TStruct) {
	if lion.LevelFatal < l.l {
		return
	}
	l.LogEntryMessage(lion.LevelFatal, newEntryMessage(event))
}

func (l *logger) Panic(event thrift.TStruct) {
	if lion.LevelPanic < l.l {
		return
	}
	l.LogEntryMessage(lion.LevelPanic, newEntryMessage(event))
}

func (l *logger) Print(event thrift.TStruct) {
	l.LogEntryMessage(lion.LevelNone, newEntryMessage(event))
}

func (l *logger) LogDebug() LevelLogger {
	if lion.LevelDebug < l.l {
		return discardLevelLoggerInstance
	}
	return newLevelLogger(l.Logger.LogDebug(), l.l)
}

func (l *logger) LogInfo() LevelLogger {
	if lion.LevelInfo < l.l {
		return discardLevelLoggerInstance
	}
	return newLevelLogger(l.Logger.LogInfo(), l.l)
}

func (l *logger) LionLogger() lion.Logger {
	return l.Logger
}

func newEntryMessage(message thrift.TStruct) *lion.EntryMessage {
	if message == nil {
		return nil
	}
	return &lion.EntryMessage{
		Encoding: Encoding,
		Value:    message,
	}
}

type levelLogger struct {
	lion.LevelLogger
	level lion.Level
}

func newLevelLogger(logger lion.LevelLogger, level lion.Level) *levelLogger {
	return &levelLogger{logger, level}
}

func (l *levelLogger) WithField(key string, value interface{}) LevelLogger {
	return &levelLogger{l.LevelLogger.WithField(key, value), l.level}
}

func (l *levelLogger) WithFields(fields map[string]interface{}) LevelLogger {
	return &levelLogger{l.LevelLogger.WithFields(fields), l.level}
}

func (l *levelLogger) WithKeyValues(keyValues ...interface{}) LevelLogger {
	return &levelLogger{l.LevelLogger.WithKeyValues(keyValues...), l.level}
}

func (l *levelLogger) WithContext(context thrift.TStruct) LevelLogger {
	return &levelLogger{l.LevelLogger.WithEntryMessageContext(newEntryMessage(context)), l.level}
}

func (l *levelLogger) Print(event thrift.TStruct) {
	l.LogEntryMessage(l.level, newEntryMessage(event))
}

type discardLevelLogger struct{}

func newDiscardLevelLogger() *discardLevelLogger {
	return &discardLevelLogger{}
}

func (d *discardLevelLogger) Printf(format string, args ...interface{})            {}
func (d *discardLevelLogger) Println(args ...interface{})                          {}
func (d *discardLevelLogger) WithField(key string, value interface{}) LevelLogger  { return d }
func (d *discardLevelLogger) WithFields(fields map[string]interface{}) LevelLogger { return d }
func (d *discardLevelLogger) WithKeyValues(keyvalues ...interface{}) LevelLogger   { return d }
func (d *discardLevelLogger) WithContext(context thrift.TStruct) LevelLogger       { return d }
func (d *discardLevelLogger) Print(event thrift.TStruct)                           {}
