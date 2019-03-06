package log

var (
	defaultLogger Logger = Build()
)

func Debug(format string, args ...interface{}) {
	defaultLogger.Debug(format, args...)
}
func DebugWith(meta MetaFields, format string, args ...interface{}) {
	defaultLogger.DebugWith(meta, format, args...)
}
func Info(format string, args ...interface{}) {
	defaultLogger.Info(format, args...)
}
func InfoWith(meta MetaFields, format string, args ...interface{}) {
	defaultLogger.InfoWith(meta, format, args...)
}
func Warning(format string, args ...interface{}) {
	defaultLogger.Warning(format, args...)
}
func WarningWith(meta MetaFields, format string, args ...interface{}) {
	defaultLogger.WarningWith(meta, format, args...)
}
func Error(format string, args ...interface{}) {
	defaultLogger.Error(format, args...)
}
func ErrorWith(meta MetaFields, format string, args ...interface{}) {
	defaultLogger.ErrorWith(meta, format, args...)
}

func Get(name string) Logger {
	return defaultLogger.Get(name)
}
