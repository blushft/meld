package server

import "log"

type Logger interface {
	Log(...interface{})
	Logf(string, ...interface{})
}

func NewLogger() Logger {
	return &logger{}
}

type logger struct{}

func (l *logger) Log(e ...interface{}) {
	log.Println(e...)
}

func (l *logger) Logf(f string, e ...interface{}) {
	log.Printf(f, e...)
}
