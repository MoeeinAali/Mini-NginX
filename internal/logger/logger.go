package logger

import (
	"log"
	"os"
)

var l *log.Logger

func Init(file string) {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	l = log.New(f, "", log.LstdFlags)
}

func Info(msg string, kv ...interface{}) {
	if l == nil {
		return
	}
	l.Println(append([]interface{}{msg}, kv...)...)
}
