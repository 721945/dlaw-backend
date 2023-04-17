package libs

import (
	"log"
)

type Logger struct {
	Logger *log.Logger
}

//var globalLogger *Logger

func NewLogger() *Logger {
	return &Logger{
		Logger: log.Default(),
	}
}

func (l Logger) GetGinLogger() *log.Logger {
	return l.Logger
}

func (l Logger) Info(message ...interface{}) {
	if len(message) > 0 {
		l.Logger.Printf("INFO: %+v --> %+v\n", message[0], message[1:])
	} else {
		l.Logger.Println("INFO: No message provided")
	}
}

func (l Logger) Error(v ...any) {
	l.Logger.Println(v)
}

func (l Logger) Fatal(v ...any) {
	l.Logger.Fatal(v)
}

func (l Logger) Panic(message string) {
	l.Logger.Panic(message)
}
