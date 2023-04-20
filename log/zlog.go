package log

type Configuration struct {
	LogLevel string
}

type LoggerInterface interface {
	Info(args ...interface{})
	InfoF(s string, args ...interface{})
	Error(args ...interface{})
	ErrorF(s string, args ...interface{})
	Warn(args ...interface{})
	WarnF(s string, args ...interface{})
}
