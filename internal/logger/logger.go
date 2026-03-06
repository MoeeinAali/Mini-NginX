package logger

import (
	"io"
	"log"
	"os"
)

var l *log.Logger

func Init(file string) {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, f)
	l = log.New(mw, "", log.LstdFlags)
}

func Info(msg string, kv ...interface{}) {
	if l == nil {
		return
	}
	l.Println(append([]interface{}{msg}, kv...)...)
}
