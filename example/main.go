package main

import (
	"context"
	"github.com/wyy-go/wlog"
	"os"
)

func main() {
	l, lvl := wlog.New(
		wlog.WithLevel("debug"),
		wlog.WithFormat("console"),
		wlog.WithStack(true),
		wlog.WithAdapter("file-custom", os.Stderr),
		wlog.WithPath(""),
		wlog.WithAddCaller(true),
		wlog.WithCallerSkip(0),

		wlog.WithFilename("log.log"),
		wlog.WithMaxSize(100),
		wlog.WithMaxAge(7),
		wlog.WithMaxBackups(7),
		wlog.WithEnableLocalTime(),
		wlog.WithEnableCompress(),
	)
	wlog.ReplaceGlobals(wlog.NewLoggerWith(l, lvl).Named("project"))
	wlog.SetDefaultValuer(
		wlog.Caller(3),
		wlog.ImmutString("field_fn_key1", "field_fn_value1"),
	)

	wlog.Debug("Debug")
	wlog.Debug("Debug", "-", "111111", "-", 1)
	wlog.Info("Info")
	wlog.Warn("Warn")
	wlog.Info("info")
	wlog.Error("Error")
	wlog.DPanic("DPanic")

	wlog.Debugf("Debugf: %s", "debug")
	wlog.Infof("Infof: %s", "info")
	wlog.Warnf("Warnf: %s", "warn")
	wlog.Infof("Infof: %s", "info")
	wlog.Errorf("Errorf: %s", "error")
	wlog.DPanicf("DPanicf: %s", "dPanic")

	wlog.Debugw("Debugw", "Debugw", "w", "11111111", "2222222222")
	wlog.Infow("Infow", "Infow", "w")
	wlog.Warnw("Warnw", "Warnw", "w")
	wlog.Infow("Infow", "Infow", "w")
	wlog.Errorw("Errorw", "Errorw", "w")
	wlog.DPanicw("DPanicw", "DPanicw", "w")

	shouPanic(func() {
		wlog.Panic("Panic")
	})
	shouPanic(func() {
		wlog.Panicf("Panicf: %s", "panic")
	})
	shouPanic(func() {
		wlog.Panicw("Panicw: %s", "panic", "w")
	})

	wlog.With(wlog.String("aa", "bb")).Debug("debug with")

	wlog.Named("another").Debug("debug named")

	wlog.WithContext(context.Background()).
		WithValuer(func(ctx context.Context) wlog.Field { return wlog.String("field_fn_key2", "field_fn_value2") }).
		Debug("with context")

	wlog.WithContext(context.Background()).
		WithValuer(func(ctx context.Context) wlog.Field { return wlog.String("field_fn_key3", "field_fn_value3") }).
		Debug("with field fn")

	wlog.With(wlog.Namespace("aaaa")).With(wlog.String("xx", "yy")).With(wlog.Namespace("bbbbbb")).With(wlog.String("dd", "gg")).Debug()

	_ = wlog.Sync()

}

func shouPanic(f func()) {
	defer func() {
		e := recover()
		if e == nil {
			wlog.Errorf("should panic but not")
		}
	}()
	f()
}
