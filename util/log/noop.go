package log

type NoOpLogger struct {
}

func (l *NoOpLogger) Debug(format string, args ...interface{}) {
}

func (l *NoOpLogger) DebugWith(meta MetaFields, format string, args ...interface{}) {
}

func (l *NoOpLogger) Info(format string, args ...interface{}) {
}

func (l *NoOpLogger) InfoWith(meta MetaFields, format string, args ...interface{}) {
}

func (l *NoOpLogger) Warning(format string, args ...interface{}) {
}

func (l *NoOpLogger) WarningWith(meta MetaFields, format string, args ...interface{}) {
}

func (l *NoOpLogger) Error(format string, args ...interface{}) {
}

func (l *NoOpLogger) ErrorWith(meta MetaFields, format string, args ...interface{}) {
}

func (l *NoOpLogger) Get(name string) Logger {
	return &NoOpLogger{}
}
