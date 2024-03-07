package utils

import "log"

type ILogger interface {
	Info(message string, data ...interface{})
	Warn(message string, data ...interface{})
	Error(message string, data ...interface{})
	Panic(err error, message string, data ...interface{})
}

type Logger struct{}

func (logger *Logger) Info(message string, data ...interface{}) {
	if len(data) > 0 {
		log.Printf(message+"; %v", data[0])
	} else {
		log.Print(message)
	}
}

func (logger *Logger) Warn(message string, data ...interface{}) {
	if len(data) > 0 {
		log.Printf(message+"; %v", data[0])
	} else {
		log.Print(message)
	}
}

func (logger *Logger) Error(message string, data ...interface{}) {
	if len(data) > 0 {
		log.Printf(message+"; %v", data[0])
	} else {
		log.Print(message)
	}
}

func (logger *Logger) Panic(err error, message string, data ...interface{}) {
	if len(data) > 0 {
		log.Panicf(message+"; %v", data[0])
	} else {
		log.Panic(message)
	}
}
