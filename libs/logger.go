package libs

import "log"

type Logger struct {
	Logger log.Logger
}

//var globalLogger *Logger

func NewLogger() Logger {
	return Logger{
		Logger: *log.Default(),
	}
}

func (l *Logger) GetGinLogger() *log.Logger {
	return &l.Logger
}

func (l *Logger) Info(message string) {
	l.Logger.Println(message)
}

func (l *Logger) Error(v ...any) {
	l.Logger.Println(v)
}

func (l *Logger) Fatal(v ...any) {
	l.Logger.Fatal(v)
}

func (l *Logger) Panic(message string) {
	l.Logger.Panic(message)
}
