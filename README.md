# wlog

[![GoDoc](https://godoc.org/github.com/wyy-go/wlog?status.svg)](https://godoc.org/github.com/wyy-go/wlog)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/wyy-go/wlog?tab=doc)
[![codecov](https://codecov.io/gh/wyy-go/wlog/branch/main/graph/badge.svg)](https://codecov.io/gh/wyy-go/wlog)
[![Go Report Card](https://goreportcard.com/badge/github.com/wyy-go/wlog)](https://goreportcard.com/report/github.com/wyy-go/wlog)
[![Licence](https://img.shields.io/github/license/wyy-go/wlog)](https://raw.githubusercontent.com/wyy-go/wlog/main/LICENSE)
[![Tag](https://img.shields.io/github/v/tag/wyy-go/wlog)](https://github.com/wyy-go/wlog/tags)

## Features


## Usage

### Installation

Use go get.
```bash
    go get github.com/wyy-go/wlog
```

Then import the package into your own code.
```bash
    import "github.com/wyy-go/wlog"
```

### Example

[embedmd]:# (examples/main.go go)
```go
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

		wlog.WithFilename("log.log"),
		wlog.WithMaxSize(100),
		wlog.WithMaxAge(7),
		wlog.WithMaxBackups(7),
		wlog.WithEnableLocalTime(),
		wlog.WithEnableCompress(),
	)
	wlog.ReplaceGlobals(wlog.NewLoggerWith(l, lvl).Named("project"))
	wlog.SetDefaultValuer(
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

```

Output:
```shell
2023-02-21 02:04:21.025	debug	project	example/main.go:30	{"msg": "Debug", "field_fn_key1": "field_fn_value1"}
2023-02-21 02:04:21.061	debug	project	example/main.go:31	{"msg": "Debug-111111-1", "field_fn_key1": "field_fn_value1"}
2023-02-21 02:04:21.062	info	project	example/main.go:32	{"msg": "Info", "field_fn_key1": "field_fn_value1"}
2023-02-21 02:04:21.062	warn	project	example/main.go:33	{"msg": "Warn", "field_fn_key1": "field_fn_value1"}
2023-02-21 02:04:21.062	info	project	example/main.go:34	{"msg": "info", "field_fn_key1": "field_fn_value1"}
2023-02-21 02:04:21.062	error	project	example/main.go:35	{"msg": "Error", "field_fn_key1": "field_fn_value1"}
2023-02-21 02:04:21.062	dpanic	project	example/main.go:36	{"msg": "DPanic", "field_fn_key1": "field_fn_value1"}
github.com/wyy-go/wlog.(*Log).DPanic
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/logger.go:182
github.com/wyy-go/wlog.DPanic
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/default.go:106
main.main
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/example/main.go:36
runtime.main
	D:/gvm/go/src/runtime/proc.go:250
2023-02-21 02:04:21.063	debug	project	example/main.go:38	{"msg": "Debugf: debug", "field_fn_key1": "field_fn_value1"}
2023-02-21 02:04:21.063	info	project	example/main.go:39	{"msg": "Infof: info", "field_fn_key1": "field_fn_value1"}
2023-02-21 02:04:21.063	warn	project	example/main.go:40	{"msg": "Warnf: warn", "field_fn_key1": "field_fn_value1"}
2023-02-21 02:04:21.063	info	project	example/main.go:41	{"msg": "Infof: info", "field_fn_key1": "field_fn_value1"}
2023-02-21 02:04:21.063	error	project	example/main.go:42	{"msg": "Errorf: error", "field_fn_key1": "field_fn_value1"}
2023-02-21 02:04:21.064	dpanic	project	example/main.go:43	{"msg": "DPanicf: dPanic", "field_fn_key1": "field_fn_value1"}
github.com/wyy-go/wlog.(*Log).DPanicf
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/logger.go:239
github.com/wyy-go/wlog.DPanicf
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/default.go:128
main.main
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/example/main.go:43
runtime.main
	D:/gvm/go/src/runtime/proc.go:250
2023-02-21 02:04:21.064	debug	project	example/main.go:45	{"msg": "Debugw", "field_fn_key1": "field_fn_value1", "Debugw": "w", "11111111": "2222222222"}
2023-02-21 02:04:21.064	info	project	example/main.go:46	{"msg": "Infow", "field_fn_key1": "field_fn_value1", "Infow": "w"}
2023-02-21 02:04:21.065	warn	project	example/main.go:47	{"msg": "Warnw", "field_fn_key1": "field_fn_value1", "Warnw": "w"}
2023-02-21 02:04:21.065	info	project	example/main.go:48	{"msg": "Infow", "field_fn_key1": "field_fn_value1", "Infow": "w"}
2023-02-21 02:04:21.065	error	project	example/main.go:49	{"msg": "Errorw", "field_fn_key1": "field_fn_value1", "Errorw": "w"}
2023-02-21 02:04:21.065	dpanic	project	example/main.go:50	{"msg": "DPanicw", "field_fn_key1": "field_fn_value1", "DPanicw": "w"}
github.com/wyy-go/wlog.(*Log).DPanicw
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/logger.go:305
github.com/wyy-go/wlog.DPanicw
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/default.go:159
main.main
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/example/main.go:50
runtime.main
	D:/gvm/go/src/runtime/proc.go:250
2023-02-21 02:04:21.065	panic	project	example/main.go:53	{"msg": "Panic", "field_fn_key1": "field_fn_value1"}
github.com/wyy-go/wlog.(*Log).Panic
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/logger.go:190
github.com/wyy-go/wlog.Panic
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/default.go:109
main.main.func1
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/example/main.go:53
main.shouPanic
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/example/main.go:87
main.main
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/example/main.go:52
runtime.main
	D:/gvm/go/src/runtime/proc.go:250
2023-02-21 02:04:21.066	panic	project	example/main.go:56	{"msg": "Panicf: panic", "field_fn_key1": "field_fn_value1"}
github.com/wyy-go/wlog.(*Log).Panicf
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/logger.go:247
github.com/wyy-go/wlog.Panicf
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/default.go:131
main.main.func2
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/example/main.go:56
main.shouPanic
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/example/main.go:87
main.main
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/example/main.go:55
runtime.main
	D:/gvm/go/src/runtime/proc.go:250
2023-02-21 02:04:21.066	panic	project	example/main.go:59	{"msg": "Panicw: %s", "field_fn_key1": "field_fn_value1", "panic": "w"}
github.com/wyy-go/wlog.(*Log).Panicw
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/logger.go:314
github.com/wyy-go/wlog.Panicw
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/default.go:163
main.main.func3
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/example/main.go:59
main.shouPanic
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/example/main.go:87
main.main
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/example/main.go:58
runtime.main
	D:/gvm/go/src/runtime/proc.go:250
2023-02-21 02:04:21.067	debug	project	example/main.go:62	{"aa": "bb", "msg": "debug with", "field_fn_key1": "field_fn_value1"}
2023-02-21 02:04:21.067	debug	project.another	example/main.go:64	{"msg": "debug named", "field_fn_key1": "field_fn_value1"}
2023-02-21 02:04:21.067	debug	project	example/main.go:68	{"msg": "with context", "field_fn_key1": "field_fn_value1", "field_fn_key2": "field_fn_value2"}
2023-02-21 02:04:21.067	debug	project	example/main.go:72	{"msg": "with field fn", "field_fn_key1": "field_fn_value1", "field_fn_key3": "field_fn_value3"}
2023-02-21 02:04:21.067	debug	project	example/main.go:74	{"aaaa": {"xx": "yy", "bbbbbb": {"dd": "gg", "msg": "", "field_fn_key1": "field_fn_value1"}}}

```