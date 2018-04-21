package mock

// Log is a mock logger implementation
type Log struct {
	InfoFn      func(fields ...interface{})
	InfoInvoked bool

	ErrorFn      func(fields ...interface{})
	ErrorInvoked bool
}

// Info level mock logs
func (l *Log) Info(fields ...interface{}) {
	if l.InfoFn == nil {
		return
	}
	l.InfoInvoked = true
	l.InfoFn(fields)
}

// Error level logs
func (l *Log) Error(fields ...interface{}) {
	if l.ErrorFn == nil {
		return
	}
	l.ErrorInvoked = true
	l.ErrorFn(fields)
}
