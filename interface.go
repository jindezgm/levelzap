/*
 * @Author: jinde.zgm
 * @Date: 2020-11-22 14:09:48
 * @Descripttion:
 */

package levelzap

import (
	"flag"

	"go.uber.org/zap/zapcore"
)

// Interface define leveled zap.Logger.
type Interface interface {
	// InitFlags is for explicitly initializing the flags.
	InitFlags(flagset *flag.FlagSet)
	// SetEncoderConfig will set the custom encoder config for zap.
	SetEncoderConfig(config zapcore.EncoderConfig)
	// AddCallerSkip increases the number of callers skipped by caller annotation
	AddCallerSkip(skip int)
	// SetLevel set number for the log level verbosity
	SetLevel(l Level)
	// V reports whether verbosity at the call site is at least the requested level.
	// The returned value is a *zap.Logger of type Verbose, which implements Info, Panic
	// and Fatal. These methods will write to the Info log if called.
	V(l Level) *Verbose
	// Flush flushes all pending log I/O.
	Flush()
}

// levelzap uses the log level defined by zap by default.
// The log level can be customized through the SetLevel() and V() interfaces.
const (
	DEBUG Level = Level(zapcore.DebugLevel)
	INFO  Level = Level(zapcore.InfoLevel)
	WARN  Level = Level(zapcore.WarnLevel)
	ERROR Level = Level(zapcore.ErrorLevel)
	PANIC Level = Level(zapcore.PanicLevel)
	FATAL Level = Level(zapcore.FatalLevel)
)

// New create levelzap object.
func New() Interface {
	l := &loggingT{}
	l.initDefault()
	return l
}

// By default, levelzap will create an Interface named logging for applications that only have one log object.
// Levelzap opens the logging's interfaces as global functions.

// InitFlags initialize logging flags.
func InitFlags(flagset *flag.FlagSet) {
	logging.InitFlags(flagset)
}

// SetEncoderConfig set logging custom encoder config.
func SetEncoderConfig(config zapcore.EncoderConfig) {
	logging.SetEncoderConfig(config)
}

// SetLevel set logging log level.
func SetLevel(l Level) {
	logging.SetLevel(l)
}

// V reports whether logging verbosity at the call site is at least the requested level.
func V(level Level) *Verbose {
	return logging.V(level)
}

// Flush flushes logging all pending log I/O.
func Flush() {
	logging.Flush()
}
