package logging

import (
	"fmt"
	"os"
	"strings"

	"github.com/bombsimon/logrusr/v2"
	"github.com/cloudogu/k8s-apply-lib/apply"
	"github.com/go-logr/logr"
	"github.com/sirupsen/logrus"
	ctrl "sigs.k8s.io/controller-runtime"
)

const logLevelEnvVar = "LOG_LEVEL"

const logWithNameFormat = "[%s] %s"

const (
	errorLevel int = iota
	warningLevel
	infoLevel
	debugLevel
)

// CurrentLogLevel is the currently configured logLevel
// The default logLevel is "ERROR"
var (
	defaultLogLevel = logrus.InfoLevel
	CurrentLogLevel = defaultLogLevel
)

type libraryLogger struct {
	logger logr.LogSink
	name   string
}

func (ll *libraryLogger) log(level int, args ...interface{}) {
	ll.logger.Info(level, fmt.Sprintf(logWithNameFormat, ll.name, fmt.Sprint(args...)))
}

func (ll *libraryLogger) logf(level int, format string, args ...interface{}) {
	ll.logger.Info(level, fmt.Sprintf(logWithNameFormat, ll.name, fmt.Sprintf(format, args...)))
}

// Debug will log a message at debug-level.
func (ll *libraryLogger) Debug(args ...interface{}) {
	ll.log(debugLevel, args...)
}

// Info will log a message at info-level.
func (ll *libraryLogger) Info(args ...interface{}) {
	ll.log(infoLevel, args...)
}

// Warning will log a message at warning-level.
func (ll *libraryLogger) Warning(args ...interface{}) {
	ll.log(warningLevel, args...)
}

// Error will log a message at error-level.
func (ll *libraryLogger) Error(args ...interface{}) {
	ll.log(errorLevel, args...)
}

// Debugf will log a message at debug-level using a format-string.
func (ll *libraryLogger) Debugf(format string, args ...interface{}) {
	ll.logf(debugLevel, format, args...)
}

// Infof will log a message at info-level using a format-string.
func (ll *libraryLogger) Infof(format string, args ...interface{}) {
	ll.logf(infoLevel, format, args...)
}

// Warningf will log a message at warning-level using a format-string.
func (ll *libraryLogger) Warningf(format string, args ...interface{}) {
	ll.logf(warningLevel, format, args...)
}

// Errorf will log a message at error-level using a format-string.
func (ll *libraryLogger) Errorf(format string, args ...interface{}) {
	ll.logf(errorLevel, format, args...)
}

func getLogLevelFromEnv() (logrus.Level, error) {
	logLevel, found := os.LookupEnv(logLevelEnvVar)
	if !found || strings.TrimSpace(logLevel) == "" {
		return defaultLogLevel, nil
	}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return defaultLogLevel, fmt.Errorf("value of log environment variable [%s] is not a valid log level: %w", logLevelEnvVar, err)
	}

	return level, nil
}

func FormattingLoggerWithName(name string, loggerFn func(msg string, keysAndValues ...interface{})) func(format string, args ...interface{}) {
	return func(format string, args ...interface{}) {
		loggerFn(fmt.Sprintf(logWithNameFormat, name, fmt.Sprintf(format, args...)))
	}
}

// ConfigureLogger configures the logger using the logLevel from the environment
func ConfigureLogger() error {
	level, err := getLogLevelFromEnv()
	if err != nil {
		return err
	}

	// create logrus logger that can be styled and formatted
	logrusLog := logrus.New()
	logrusLog.SetFormatter(&logrus.TextFormatter{})
	logrusLog.SetLevel(level)

	CurrentLogLevel = level

	// convert logrus logger to logr logger
	logrusLogrLogger := logrusr.New(logrusLog)

	// set logr logger as controller logger
	ctrl.SetLogger(logrusLogrLogger)

	// set custom logger implementation to k8s-apply-lib logger
	k8sApplyLibLogger := libraryLogger{name: "k8s-apply-lib", logger: logrusLogrLogger.GetSink()}
	apply.GetLogger = func() apply.Logger {
		return &k8sApplyLibLogger
	}

	return nil
}
