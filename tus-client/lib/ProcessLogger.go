package lib

import (
	"io"
	"log"
	"os"
)

func LogInit(l LogUtils) {
	l.ProcessLogger()
}

type LogUtils interface {
	ProcessLogger()
}

type Logs struct {
	LogPath string
}

func (l *Logs) ProcessLogger() {
	fpLog, err := os.OpenFile(l.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	defer fpLog.Close()

	multiWriter := io.MultiWriter(fpLog, os.Stdout)

	log.SetOutput(multiWriter)
	log.SetFlags(0)
}
