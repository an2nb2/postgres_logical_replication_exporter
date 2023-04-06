package main

import (
	"fmt"

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

func newPromLogger(logger log.Logger) *promLogger {
	return &promLogger{logger}
}

func newLogger(level string) (log.Logger, error) {
	loglevel := &promlog.AllowedLevel{}
	err := loglevel.Set(level)

	if err != nil {
		return nil, err
	}

	config := &promlog.Config{
		Level: loglevel,
	}

	return promlog.New(config), nil
}
