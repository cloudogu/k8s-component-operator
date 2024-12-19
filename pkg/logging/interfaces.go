package logging

import "github.com/go-logr/logr"

type logSink interface {
	logr.LogSink
}
