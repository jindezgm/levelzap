/*
 * @Author: jinde.zgm
 * @Date: 2020-11-22 14:22:48
 * @Descripttion:
 */

package levelzap

import (
	"flag"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// LogFileDefaultMaxSizeMB define log file default max size in MB.
	LogFileDefaultMaxSizeMB = 1024
	// LogFileDefaultMaxBackups define log file default max backups.
	LogFileDefaultMaxBackups = 7
	// LogFileDefaultMaxAge define log file default max age.
	LogFileDefaultMaxAge = 7
)

// init sets up the defaults and runs flushDaemon.
func init() {
	logging.initDefault()
	logging.AddCallerSkip(1)
}

// loggingT collects all the global state of the logging setup.
type loggingT struct {
	// Boolean flags. Not handled atomically because the flag.Value interface
	// does not let us avoid the =true, and that shorthand is necessary for
	// compatibility.
	toStderr     bool // The --logtostderr flag.
	alsoToStderr bool // The --logalsotostderr flag.
	addCaller    bool // The --logaddcaller flag.

	// mutex protects the remaining elements of this structure and is used to synchronize logging.
	mutex sync.Mutex

	// V logging level, the value of the -v flag
	verbosity Level
	// If non-empty, overrides the choice of directory in which to write logs.
	// See createLogDirs for the full list of possible destinations.
	logDir string

	// If non-empty, specifies the path of the file to write logs. mutually exclusive
	// with the log_dir option.
	logFile string

	// When logFile is specified, this limiter makes sure the logFile won't exceeds a certain size. When exceeds, the
	// logFile will be cleaned up. If this value is 0, no size limitation will be applied to logFile.
	logFileMaxSizeMB int

	// logFileMaxBackups is the maximum number of old log files to retain.  The default
	// is to retain all old log files (though MaxAge may still cause them to get deleted.)
	logFileMaxBackups int
	// logfileMaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	logFileMaxAge int

	// callerSkip define the number of callers skipped by caller annotation.
	callerSkip int

	// If set, all output will be redirected unconditionally to the provided logr.Logger
	logger        atomic.Value
	encoderConfig *zapcore.EncoderConfig
}

var logging loggingT

// initDefault initialize logging default configure.
func (l *loggingT) initDefault() {
	l.addCaller = true
	l.callerSkip = 1
	l.logFileMaxSizeMB = LogFileDefaultMaxSizeMB
	l.logFileMaxBackups = LogFileDefaultMaxBackups
	l.logFileMaxAge = LogFileDefaultMaxAge
	l.verbosity = INFO
	l.logFile = filepath.Base(os.Args[0]) + ".log"
	l.logDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
}

// getLogger create logger if nil.
func (l *loggingT) getLogger() *zap.Logger {
	// Check logger nil or not.
	logger := l.logger.Load()
	if nil == logger {
		l.mutex.Lock()
		// Check logger nil or not again.
		if logger = l.logger.Load(); nil == logger {
			if nil == l.encoderConfig {
				l.encoderConfig = new(zapcore.EncoderConfig)
				*l.encoderConfig = zap.NewProductionEncoderConfig()
				l.encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
			}
			// Create zap core.
			var core zapcore.Core
			if l.toStderr {
				// stderr only.
				core = zapcore.NewCore(zapcore.NewConsoleEncoder(*l.encoderConfig), zapcore.AddSync(os.Stderr), zap.InfoLevel)
			} else {
				// Create file writer.
				file := zapcore.AddSync(&lumberjack.Logger{
					Filename:   filepath.Join(l.logDir, l.logFile),
					MaxSize:    l.logFileMaxSizeMB,
					MaxBackups: l.logFileMaxBackups,
					MaxAge:     l.logFileMaxAge,
				})
				if !l.alsoToStderr {
					// File only.
					core = zapcore.NewCore(zapcore.NewConsoleEncoder(*l.encoderConfig), file, zap.InfoLevel)
				} else {
					// stderr and file.
					core = zapcore.NewTee(
						zapcore.NewCore(zapcore.NewConsoleEncoder(*l.encoderConfig), file, zap.InfoLevel),
						zapcore.NewCore(zapcore.NewConsoleEncoder(*l.encoderConfig), zapcore.AddSync(os.Stderr), zap.InfoLevel),
					)
				}
			}
			// Create zap logger.
			if l.addCaller {
				logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(l.callerSkip), zap.AddStacktrace(zap.PanicLevel))
			} else {
				logger = zap.New(core, zap.AddStacktrace(zap.PanicLevel))
			}
			l.logger.Store(logger)
		}
		l.mutex.Unlock()
	}
	return logger.(*zap.Logger)
}

// loggingT implement levelzap.Interface.var _ Interface = &loggingT{}
// InitFlags loggingT implement Interface.InitFlags.
func (l *loggingT) InitFlags(flagset *flag.FlagSet) {
	if flagset == nil {
		flagset = flag.CommandLine
	}
	flagset.BoolVar(&l.toStderr, "logtostderr", l.toStderr, "log to standard error instead of files")
	flagset.BoolVar(&l.alsoToStderr, "logalsotostderr", l.alsoToStderr, "log to standard error as well as files")
	flagset.BoolVar(&l.addCaller, "logaddcaller", l.addCaller, "annotate each message with the filename, line number, and function name of levelzap's caller")
	flagset.StringVar(&l.logDir, "logdir", l.logDir, "If non-empty, write log files in this directory")
	flagset.StringVar(&l.logFile, "logfile", l.logFile, "If non-empty, use this log file")
	flagset.IntVar(&l.logFileMaxSizeMB, "logfilemaxsize", l.logFileMaxSizeMB, "Defines the maximum size a log file can grow to. Unit is megabytes. If the value is 0, the maximum file size is unlimited.")
	flagset.IntVar(&l.logFileMaxBackups, "logfilemaxbackups", l.logFileMaxBackups, "Defines the maximum number of old log files to retain")
	flagset.IntVar(&l.logFileMaxAge, "logfilemaxage", l.logFileMaxAge, "Defines he maximum number of days to retain old log files based on the timestamp encoded in their filename.")
	flagset.Var(&l.verbosity, "v", "number for the log level verbosity")
}

// AddCallerSkip loggingT implement Interface.AddCallerSkip.
func (l *loggingT) AddCallerSkip(skip int) { l.callerSkip += skip }

// SetEncoderConfig loggingT implement Interface.SetEncoderConfig.
func (l *loggingT) SetEncoderConfig(config zapcore.EncoderConfig) {
	l.encoderConfig = new(zapcore.EncoderConfig)
	*l.encoderConfig = config
}

// SetLevel loggingT implement Interface.SetLevel.
func (l *loggingT) SetLevel(level Level) { l.verbosity.set(level) }

// V loggingT implement Interface.V.
func (l *loggingT) V(level Level) *Verbose {
	// Here is a cheap but safe test to see if V logging is enabled globally.
	if l.verbosity.get() <= level {
		return (*Verbose)(l.getLogger())
	}
	return nil
}

// Flush loggingT implement Interface.Flush.
func (l *loggingT) Flush() { l.getLogger().Sync() }
