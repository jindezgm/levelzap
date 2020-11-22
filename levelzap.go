/*
 * @Author: jinde.zgm
 * @Date: 2020-11-22 14:09:09
 * @Descripttion:
 */

package levelzap

import (
	"strconv"
	"sync/atomic"

	"go.uber.org/zap"
)

// Level specifies a level of verbosity for V logs. *Level implements flag.Value;
// the -v flag is of type Level and should be modified only through the flag.Value interface.
type Level int32

// get returns the value of the Level.
func (l *Level) get() Level {
	return Level(atomic.LoadInt32((*int32)(l)))
}

// set sets the value of the Level.
func (l *Level) set(val Level) {
	atomic.StoreInt32((*int32)(l), int32(val))
}

// String is part of the flag.Value interface.
func (l *Level) String() string {
	return strconv.FormatInt(int64(*l), 10)
}

// Get is part of the flag.Getter interface.
func (l *Level) Get() interface{} {
	return *l
}

// Set is part of the flag.Value interface.
func (l *Level) Set(value string) error {
	val, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	atomic.StoreInt32((*int32)(l), int32(val))
	return nil
}

// Verbose is a boolean type that implements Info (like Printf) etc.
// See the documentation of V for more information.
type Verbose zap.Logger

// Info is equivalent to the zap.Logger.Info function, guarded by the value of v.
func (v *Verbose) Info(msg string, fields ...zap.Field) {
	if nil != v {
		(*zap.Logger)(v).Info(msg, fields...)
	}
}

// Panic is equivalent to the zap.Logger.Panic function, guarded by the value of v.
func (v *Verbose) Panic(msg string, fields ...zap.Field) {
	if nil != v {
		(*zap.Logger)(v).Panic(msg, fields...)
	}
}

// Fatal is equivalent to the zap.Logger.Fatal function, guarded by the value of v.
func (v *Verbose) Fatal(msg string, fields ...zap.Field) {
	if nil != v {
		(*zap.Logger)(v).Fatal(msg, fields...)
	}
}

// Debug write DEBUG level log
func Debug(msg string, fields ...zap.Field) { V(DEBUG).Info(msg, fields...) }

// Info write INFO level log
func Info(msg string, fields ...zap.Field) { V(INFO).Info(msg, fields...) }

// Warn write WARN level log
func Warn(msg string, fields ...zap.Field) { V(WARN).Info(msg, fields...) }

// Error write ERROR level log
func Error(msg string, fields ...zap.Field) { V(ERROR).Info(msg, fields...) }

// Panic write PANIC level log
func Panic(msg string, fields ...zap.Field) { V(PANIC).Panic(msg, fields...) }

// Fatal write FATAL level log
func Fatal(msg string, fields ...zap.Field) { V(FATAL).Fatal(msg, fields...) }
