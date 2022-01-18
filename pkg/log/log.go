package log

// InfoLogger represents the ability to log non-error messages, at a particular verbosity.
type InfoLogger interface {
	Info(msg string, fields ...Field)
	Infof(format string, v ...interface{})
	Infow(msg string, keysAndValues ...interface{})
}
