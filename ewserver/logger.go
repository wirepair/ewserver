package ewserver

// LogService handles logging either info or errors
type LogService interface {
	Info(fields ...interface{})
	Error(fields ...interface{})
}
