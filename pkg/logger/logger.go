package logger

import (
	"go.uber.org/zap"
)

type Config struct {
	EnableConsole     bool
	ConsoleWriteAsync bool
	ConsoleJSONFormat bool
	ConsoleLevel      string
	EnableFile        bool
	FileWriteAsync    bool
	FileJSONFormat    bool
	FileLevel         string
	FileName          string
	MaxSize           int
	MaxBackups        int
	Compress          bool
	AddCaller         bool
	MaxFileHistory    int
}

const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	FatalLevel = "fatal"
)

type Logger struct {
	raw          *zap.Logger
	sugared      *zap.SugaredLogger
	consoleLevel zap.AtomicLevel
	fileLevel    zap.AtomicLevel
}

type Fields map[string]interface{}

var (
	_defaultConfig = Config{
		EnableConsole:     true,
		ConsoleWriteAsync: false,
		ConsoleJSONFormat: false,
		ConsoleLevel:      DebugLevel,
		EnableFile:        false,
		FileWriteAsync:    false,
		FileJSONFormat:    true,
		FileLevel:         DebugLevel,
		FileName:          "go-one",
		MaxSize:           20,
		MaxBackups:        100,
		Compress:          true,
		AddCaller:         true,
		MaxFileHistory:    30,
	}
)

var _log *Logger

func Debugf(format string, args ...interface{}) {
	_log.sugared.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	_log.sugared.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	_log.sugared.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	_log.sugared.Errorf(format, args...)
}

func DPanicf(format string, args ...interface{}) {
	_log.sugared.DPanicf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	_log.sugared.Panicf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	_log.sugared.Fatalf(format, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	_log.sugared.Debugw(msg, keysAndValues...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	_log.sugared.Infow(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	_log.sugared.Warnw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	_log.sugared.Errorw(msg, keysAndValues...)
}

func DPanicw(msg string, keysAndValues ...interface{}) {
	_log.sugared.DPanicw(msg, keysAndValues...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	_log.sugared.Panicw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	_log.sugared.Fatalw(msg, keysAndValues...)
}

func WithFields(fields Fields) *Logger {
	var f = make([]interface{}, 0, len(fields))
	for k, v := range fields {
		f = append(f, k)
		f = append(f, v)
	}

	sugar := _log.raw.WithOptions(zap.AddCallerSkip(1)).Sugar().With(f...)
	return &Logger{
		sugared: sugar,
	}
}

func Zap() *zap.Logger {
	return _log.raw
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.sugared.Debugf(format, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.sugared.Infof(format, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.sugared.Warnf(format, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.sugared.Errorf(format, args...)
}

func (l *Logger) DPanicf(format string, args ...interface{}) {
	l.sugared.DPanicf(format, args...)
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	l.sugared.Panicf(format, args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.sugared.Fatalf(format, args...)
}

func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.sugared.Debugw(msg, keysAndValues...)
}

func (l *Logger) Infow(msg string, keysAndValues ...interface{}) {
	l.sugared.Infow(msg, keysAndValues...)
}

func (l *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	l.sugared.Warnw(msg, keysAndValues...)
}

func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	l.sugared.Errorw(msg, keysAndValues...)
}

func (l *Logger) DPanicw(msg string, keysAndValues ...interface{}) {
	l.sugared.DPanicw(msg, keysAndValues...)
}

func (l *Logger) Panicw(msg string, keysAndValues ...interface{}) {
	l.sugared.Panicw(msg, keysAndValues...)
}

func (l *Logger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.sugared.Fatalw(msg, keysAndValues...)
}

func (l *Logger) Zap() *zap.Logger {
	return l.raw
}

func Init(config *Config) {
	logger, err := newZapLogger(config)
	if err != nil {
		panic(err)
	}
	_log = logger
}

func NewLogger(config *Config) *Logger {
	l, err := newZapLogger(config)
	if err != nil {
		panic(err)
	}
	return l
}

func init() {
	logger, err := newZapLogger(&_defaultConfig)
	if err != nil {
		panic(err)
	}
	_log = logger
}
