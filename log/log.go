package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Info(message string)
	Warn(message string)
	Error(message string, err error)
	WithFields(fields Fields)
	InfoWithFields(message string, fields Fields)
	WarnWithFields(message string, fields Fields)
	ErrorWithFields(message string, err error, fields Fields)
}

type logger struct {
	logger *logrus.Entry
}

type Fields map[string]interface{}

var (
	WarnLevel  = logrus.WarnLevel
	ErrorLevel = logrus.ErrorLevel
	FatalLevel = logrus.FatalLevel
	InfoLevel  = logrus.InfoLevel
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	logrus.SetOutput(os.Stdout)

	logrus.SetLevel(logrus.InfoLevel)
}

func NewLogger() *logger {
	return &logger{
		logger: logrus.WithFields(logrus.Fields{}),
	}
}

func (logger *logger) AddHook(hook logrus.Hook) {
	logrus.AddHook(hook)
}

func convertToLogrusFields(fields Fields) logrus.Fields {
	logrusFields := logrus.Fields{}

	for k, v := range fields {
		logrusFields[k] = v
	}

	return logrusFields
}

// WithFields adds the specified fields to the logger permanently. All future calls to Info, Warn, Error or Fatal will
// include these fields.
func (logger *logger) WithFields(fields Fields) {
	logger.logger = logger.logger.WithFields(convertToLogrusFields(fields))
}

func (logger *logger) Info(message string) {
	logger.logger.Info(message)
}

func (logger *logger) InfoWithFields(message string, fields Fields) {
	logger.logger.WithFields(convertToLogrusFields(fields)).Info(message)
}

func (logger *logger) Warn(message string) {
	logger.logger.Warn(message)
}

func (logger *logger) WarnWithFields(message string, fields Fields) {
	logger.logger.WithFields(convertToLogrusFields(fields)).Warn(message)
}

func (logger *logger) Error(message string, err error) {
	msg := message + ": " + err.Error()
	logger.logger.Error(msg)
}

func (logger *logger) ErrorWithFields(message string, err error, fields Fields) {
	msg := message + ": " + err.Error()
	logger.logger.WithFields(convertToLogrusFields(fields)).Error(msg)
}
