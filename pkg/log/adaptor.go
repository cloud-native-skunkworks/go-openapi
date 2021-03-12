package logadaptor

import log "github.com/sirupsen/logrus"

// This adaptor is required to work with Jaeger logging
type LogrusAdapter struct {
	Logger *log.Logger
}

func (l LogrusAdapter) Error(msg string) {
	l.Logger.Errorf(msg)
}

func (l LogrusAdapter) Infof(msg string, args ...interface{}) {
	l.Logger.Infof(msg, args...)
}
