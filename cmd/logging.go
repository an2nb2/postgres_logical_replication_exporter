package main

import (
	"fmt"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/common/promlog"
)

type promLogger struct {
	logger log.Logger
}

func (l promLogger) Println(v ...interface{}) {
	_ = level.Error(l.logger).Log(fmt.Sprint(v...))
}

// Returns a promLogger instance which implements promhttp.Logger interface.
func newPromLogger(logger log.Logger) *promLogger {
	return &promLogger{logger}
}

// New returns a new go-kit logger with log level taken from arguments.
// Each logged line will be annotated with a timestamp.
func newLogger(level string) log.Logger {
	loglevel := &promlog.AllowedLevel{}
	err := loglevel.Set(level)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to initialize logger: %v\n", err)
		os.Exit(1)
	}

	config := &promlog.Config{
		Level: loglevel,
	}

	return promlog.New(config)
}
