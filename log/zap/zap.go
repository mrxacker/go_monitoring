package zap

import (
	"fmt"
	"monitoring/log"
	"net"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// var _ log.Logger = &Logger{}

type Logger struct {
	L *zap.Logger
	c net.Conn
}

type LoggerFile struct {
	L *zap.Logger
}

func NewLogger(addr string) (*Logger, func(), error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, nil, err
	}

	cfg := zap.NewProductionEncoderConfig()
	cfg.TimeKey = "timestamp"
	cfg.LevelKey = "severity"
	cfg.MessageKey = "message"

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg),
		zapcore.AddSync(conn),
		zapcore.InfoLevel,
	)

	z := zap.New(core).With(
		zap.String("service", "logger-practice"),
		zap.String("env", "dev"))

	l := &Logger{
		L: z,
		c: conn,
	}

	cleanup := func() {
		_ = l.L.Sync()
		_ = l.c.Close()
	}

	return l, cleanup, nil
}

func NewLoggerFile(filePath string, serviceName string) (*LoggerFile, func(), error) {
	// Open or create log file for appending
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, err
	}

	cfg := zap.NewProductionEncoderConfig()
	cfg.TimeKey = "timestamp"
	cfg.LevelKey = "severity"
	cfg.MessageKey = "message"

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg),
		zapcore.AddSync(f),
		zapcore.InfoLevel,
	)

	z := zap.New(core).With(
		zap.String("service", serviceName),
		zap.String("env", "dev"),
	)

	l := &LoggerFile{
		L: z,
	}

	cleanup := func() {
		_ = l.L.Sync()
		_ = f.Close()
	}

	return l, cleanup, nil
}

func (l *LoggerFile) Trace(msg string, fields ...log.Field) {
	if ce := l.L.Check(zap.DebugLevel, msg); ce != nil {
		ce.Write(zapifyFields(fields...)...)
	}
}

// Tracef logs at Trace log level using fmt formatter
func (l *LoggerFile) Tracef(msg string, args ...interface{}) {
	if ce := l.L.Check(zap.DebugLevel, ""); ce != nil {
		ce.Message = fmt.Sprintf(msg, args...)
		ce.Write()
	}
}

// Debug logs at Debug log level using fields
func (l *LoggerFile) Debug(msg string, fields ...log.Field) {
	if ce := l.L.Check(zap.DebugLevel, msg); ce != nil {
		ce.Write(zapifyFields(fields...)...)
	}
}

// Debugf logs at Debug log level using fmt formatter
func (l *LoggerFile) Debugf(msg string, args ...interface{}) {
	if ce := l.L.Check(zap.DebugLevel, ""); ce != nil {
		ce.Message = fmt.Sprintf(msg, args...)
		ce.Write()
	}
}

// Info logs at Info log level using fields
func (l *LoggerFile) Info(msg string, fields ...log.Field) {
	if ce := l.L.Check(zap.InfoLevel, msg); ce != nil {
		ce.Write(zapifyFields(fields...)...)
	}
}

// Infof logs at Info log level using fmt formatter
func (l *LoggerFile) Infof(msg string, args ...interface{}) {
	if ce := l.L.Check(zap.InfoLevel, ""); ce != nil {
		ce.Message = fmt.Sprintf(msg, args...)
		ce.Write()
	}
}

// Warn logs at Warn log level using fields
func (l *LoggerFile) Warn(msg string, fields ...log.Field) {
	if ce := l.L.Check(zap.WarnLevel, msg); ce != nil {
		ce.Write(zapifyFields(fields...)...)
	}
}

// Warnf logs at Warn log level using fmt formatter
func (l *LoggerFile) Warnf(msg string, args ...interface{}) {
	if ce := l.L.Check(zap.WarnLevel, ""); ce != nil {
		ce.Message = fmt.Sprintf(msg, args...)
		ce.Write()
	}
}

// Error logs at Error log level using fields
func (l *LoggerFile) Error(msg string, fields ...log.Field) {
	if ce := l.L.Check(zap.ErrorLevel, msg); ce != nil {
		ce.Write(zapifyFields(fields...)...)
	}
}

// Errorf logs at Error log level using fmt formatter
func (l *LoggerFile) Errorf(msg string, args ...interface{}) {
	if ce := l.L.Check(zap.ErrorLevel, ""); ce != nil {
		ce.Message = fmt.Sprintf(msg, args...)
		ce.Write()
	}
}

// Fatal logs at Fatal log level using fields
func (l *LoggerFile) Fatal(msg string, fields ...log.Field) {
	if ce := l.L.Check(zap.FatalLevel, msg); ce != nil {
		ce.Write(zapifyFields(fields...)...)
	}
}

// Fatalf logs at Fatal log level using fmt formatter
func (l *LoggerFile) Fatalf(msg string, args ...interface{}) {
	if ce := l.L.Check(zap.FatalLevel, ""); ce != nil {
		ce.Message = fmt.Sprintf(msg, args...)
		ce.Write()
	}
}
