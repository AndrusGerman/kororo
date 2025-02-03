package ports

type LogService interface {
	Info(module string, message string)
	Error(module string, message string)
}
