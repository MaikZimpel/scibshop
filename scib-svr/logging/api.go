package logging

import (
	"golang.org/x/net/context"
)

type (
	Logger interface {
		Debug(c context.Context, format string, args ...interface{})
		Info(c context.Context, format string, args ...interface{})
		Warning(c context.Context, format string, args ...interface{})
		Error(c context.Context, format string, args ...interface{})
		Critical(c context.Context, format string, args ...interface{})
	}

	loggerFactory func() Logger
)

var (
	New loggerFactory
)
