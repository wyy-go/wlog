package wlog

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Valuer is returns a log value.
type Valuer func(ctx context.Context) Field

var defaultMessageKey = "msg"

// Log wrap zap logger
type Log struct {
	log   *zap.Logger
	level zap.AtomicLevel
	fn    []Valuer
	ctx   context.Context
}

// NewLoggerWith new logger with zap logger and atomic level
func NewLoggerWith(logger *zap.Logger, lv zap.AtomicLevel) *Log {
	return &Log{
		log:   logger,
		level: lv,
		fn:    nil,
		ctx:   context.Background(),
	}
}

// NewLogger new logger
func NewLogger(opts ...Option) *Log { return NewLoggerWith(New(opts...)) }

// SetLevelWithText alters the logging level.
// ParseAtomicLevel set the logging level based on a lowercase or all-caps ASCII
// representation of the log level.
// If the provided ASCII representation is
// invalid an error is returned.
// see zapcore.Level
func (l *Log) SetLevelWithText(text string) error {
	lv, err := zapcore.ParseLevel(text)
	if err != nil {
		return err
	}
	l.level.SetLevel(lv)
	return nil
}

// SetLevel alters the logging level.
func (l *Log) SetLevel(lv zapcore.Level) *Log {
	l.level.SetLevel(lv)
	return l
}

// SetDefaultValuer set default Valuer function, which hold always until you call WithContext.
func (l *Log) SetDefaultValuer(fs ...Valuer) *Log {
	fn := make([]Valuer, 0, len(fs)+len(l.fn))
	fn = append(fn, l.fn...)
	fn = append(fn, fs...)
	l.fn = fn
	return l
}

// GetLevel returns the minimum enabled log level.
func (l *Log) GetLevel() zapcore.Level { return l.level.Level() }

// Enabled returns true if the given level is at or above this level.
func (l *Log) Enabled(lvl zapcore.Level) bool { return l.level.Enabled(lvl) }

// V returns true if the given level is at or above this level.
// same as Enabled
func (l *Log) V(lvl int) bool { return l.level.Enabled(zapcore.Level(lvl)) }

// Sugar wraps the Logger to provide a more ergonomic, but slightly slower,
// API. Sugaring a Logger is quite inexpensive, so it's reasonable for a
// single application to use both Loggers and SugaredLoggers, converting
// between them on the boundaries of performance-sensitive code.
func (l *Log) Sugar() *zap.SugaredLogger { return l.log.Sugar() }

// Logger return internal logger
func (l *Log) Logger() *zap.Logger { return l.log }

// WithValuer with Valuer function.
func (l *Log) WithValuer(fs ...Valuer) *Log {
	fn := make([]Valuer, 0, len(fs)+len(l.fn))
	fn = append(fn, l.fn...)
	fn = append(fn, fs...)
	return &Log{
		log:   l.log,
		level: l.level,
		fn:    fn,
		ctx:   l.ctx,
	}
}

// WithNewValuer return log with new Valuer function without default Valuer.
func (l *Log) WithNewValuer(fs ...Valuer) *Log {
	return &Log{
		log:   l.log,
		level: l.level,
		fn:    fs,
		ctx:   l.ctx,
	}
}

// WithContext return log with inject context.
func (l *Log) WithContext(ctx context.Context) *Log {
	return &Log{
		log:   l.log,
		level: l.level,
		fn:    l.fn,
		ctx:   ctx,
	}
}

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa.
func (l *Log) With(fields ...Field) *Log {
	return &Log{
		log:   l.log.With(fields...),
		level: l.level,
		fn:    l.fn,
		ctx:   l.ctx,
	}
}

// Named adds a sub-scope to the logger's name. See Log.Named for details.
func (l *Log) Named(name string) *Log {
	return &Log{
		log:   l.log.Named(name),
		level: l.level,
		fn:    l.fn,
		ctx:   l.ctx,
	}
}

// Sync flushes any buffered log entries.
func (l *Log) Sync() error {
	return l.log.Sync()
}

// Debug uses fmt.Sprint to construct and log a message.
func (l *Log) Debug(args ...any) {
	if !l.level.Enabled(DebugLevel) {
		return
	}
	l.log.With(getMessage("", args)).With(injectField(l.ctx, l.fn)...).Sugar().Debug()
}

// Info uses fmt.Sprint to construct and log a message.
func (l *Log) Info(args ...any) {
	if !l.level.Enabled(InfoLevel) {
		return
	}
	l.log.With(getMessage("", args)).With(injectField(l.ctx, l.fn)...).Sugar().Info()
}

// Warn uses fmt.Sprint to construct and log a message.
func (l *Log) Warn(args ...any) {
	if !l.level.Enabled(WarnLevel) {
		return
	}
	l.log.With(getMessage("", args)).With(injectField(l.ctx, l.fn)...).Sugar().Warn()
}

// Error uses fmt.Sprint to construct and log a message.
func (l *Log) Error(args ...any) {
	if !l.level.Enabled(ErrorLevel) {
		return
	}
	l.log.With(getMessage("", args)).With(injectField(l.ctx, l.fn)...).Sugar().Error()
}

// DPanic uses fmt.Sprint to construct and log a message. In development, the
// logger then panics. (See DPanicLevel for details.)
func (l *Log) DPanic(args ...any) {
	if !l.level.Enabled(DPanicLevel) {
		return
	}
	l.log.With(getMessage("", args)).With(injectField(l.ctx, l.fn)...).Sugar().DPanic()
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func (l *Log) Panic(args ...any) {
	if !l.level.Enabled(PanicLevel) {
		return
	}
	l.log.With(getMessage("", args)).With(injectField(l.ctx, l.fn)...).Sugar().Panic()
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func (l *Log) Fatal(args ...any) {
	if !l.level.Enabled(FatalLevel) {
		return
	}
	l.log.With(getMessage("", args)).With(injectField(l.ctx, l.fn)...).Sugar().Fatal()
}

// Debugf uses fmt.Sprintf to log a templated message.
func (l *Log) Debugf(template string, args ...any) {
	if !l.level.Enabled(DebugLevel) {
		return
	}
	l.log.With(getMessage(template, args)).With(injectField(l.ctx, l.fn)...).Sugar().Debugf(template, args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func (l *Log) Infof(template string, args ...any) {
	if !l.level.Enabled(InfoLevel) {
		return
	}
	l.log.With(getMessage(template, args)).With(injectField(l.ctx, l.fn)...).Sugar().Infof(template, args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func (l *Log) Warnf(template string, args ...any) {
	if !l.level.Enabled(WarnLevel) {
		return
	}
	l.log.With(getMessage(template, args)).With(injectField(l.ctx, l.fn)...).Sugar().Warnf(template, args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func (l *Log) Errorf(template string, args ...any) {
	if !l.level.Enabled(ErrorLevel) {
		return
	}
	l.log.With(getMessage(template, args)).With(injectField(l.ctx, l.fn)...).Sugar().Errorf(template, args...)
}

// DPanicf uses fmt.Sprintf to log a templated message. In development, the
// logger then panics. (See DPanicLevel for details.)
func (l *Log) DPanicf(template string, args ...any) {
	if !l.level.Enabled(DPanicLevel) {
		return
	}
	l.log.With(getMessage(template, args)).With(injectField(l.ctx, l.fn)...).Sugar().DPanicf(template, args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func (l *Log) Panicf(template string, args ...any) {
	if !l.level.Enabled(PanicLevel) {
		return
	}
	l.log.With(getMessage(template, args)).With(injectField(l.ctx, l.fn)...).Sugar().Panicf(template, args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func (l *Log) Fatalf(template string, args ...any) {
	if !l.level.Enabled(FatalLevel) {
		return
	}
	l.log.With(getMessage(template, args)).With(injectField(l.ctx, l.fn)...).Sugar().Fatalf(template, args...)
}

// Debugw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
//
// When debug-level logging is disabled, this is much faster than
//
//	s.With(keysAndValues).Debug(msg)
func (l *Log) Debugw(msg string, keysAndValues ...any) {
	if !l.level.Enabled(DebugLevel) {
		return
	}
	l.log.With(getMessage(msg, nil)).With(injectField(l.ctx, l.fn)...).Sugar().Debugw("", keysAndValues...)
}

// Infow logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (l *Log) Infow(msg string, keysAndValues ...any) {
	if !l.level.Enabled(InfoLevel) {
		return
	}
	l.log.With(getMessage(msg, nil)).With(injectField(l.ctx, l.fn)...).Sugar().Infow("", keysAndValues...)
}

// Warnw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (l *Log) Warnw(msg string, keysAndValues ...any) {
	if !l.level.Enabled(WarnLevel) {
		return
	}
	l.log.With(getMessage(msg, nil)).With(injectField(l.ctx, l.fn)...).Sugar().Warnw("", keysAndValues...)
}

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func (l *Log) Errorw(msg string, keysAndValues ...any) {
	if !l.level.Enabled(ErrorLevel) {
		return
	}
	l.log.With(getMessage(msg, nil)).With(injectField(l.ctx, l.fn)...).Sugar().Errorw("", keysAndValues...)
}

// DPanicw logs a message with some additional context. In development, the
// logger then panics. (See DPanicLevel for details.) The variadic key-value
// pairs are treated as they are in With.
func (l *Log) DPanicw(msg string, keysAndValues ...any) {
	if !l.level.Enabled(DPanicLevel) {
		return
	}
	l.log.With(getMessage(msg, nil)).With(injectField(l.ctx, l.fn)...).Sugar().DPanicw("", keysAndValues...)
}

// Panicw logs a message with some additional context, then panics. The
// variadic key-value pairs are treated as they are in With.
func (l *Log) Panicw(msg string, keysAndValues ...any) {
	if !l.level.Enabled(PanicLevel) {
		return
	}
	l.log.With(getMessage(msg, nil)).With(injectField(l.ctx, l.fn)...).Sugar().Panicw("", keysAndValues...)
}

// Fatalw logs a message with some additional context, then calls os.Exit. The
// variadic key-value pairs are treated as they are in With.
func (l *Log) Fatalw(msg string, keysAndValues ...any) {
	if !l.level.Enabled(FatalLevel) {
		return
	}
	l.log.With(getMessage(msg, nil)).With(injectField(l.ctx, l.fn)...).Sugar().Fatalw("", keysAndValues...)
}

func injectField(ctx context.Context, vs []Valuer) []Field {
	var fields []Field

	if len(vs) > 0 {
		fields = make([]Field, 0, len(vs))
		for _, f := range vs {
			fields = append(fields, f(ctx))
		}
	}
	return fields
}

func getMessage(template string, fmtArgs []interface{}) zap.Field {
	if len(fmtArgs) == 0 {
		return zap.String(defaultMessageKey, template)
	}

	if template != "" {
		return zap.String(defaultMessageKey, fmt.Sprintf(template, fmtArgs...))
	}

	if len(fmtArgs) == 1 {
		if str, ok := fmtArgs[0].(string); ok {
			return zap.String(defaultMessageKey, str)
		}
	}
	return zap.String(defaultMessageKey, fmt.Sprint(fmtArgs...))
}
