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


```

Output:
```shell
2023-02-21 01:11:08.924	debug	project	wlog/logger.go:149	{"msg": "Debug", "caller": "main.go:32", "field_fn_key1": "field_fn_value1"}
2023-02-21 01:11:08.971	debug	project	wlog/logger.go:149	{"msg": "Debug-111111-1", "caller": "main.go:33", "field_fn_key1": "field_fn_value1"}
2023-02-21 01:11:08.971	info	project	wlog/logger.go:157	{"msg": "Info", "caller": "main.go:34", "field_fn_key1": "field_fn_value1"}
2023-02-21 01:11:08.972	warn	project	wlog/logger.go:165	{"msg": "Warn", "caller": "main.go:35", "field_fn_key1": "field_fn_value1"}
2023-02-21 01:11:08.972	info	project	wlog/logger.go:157	{"msg": "info", "caller": "main.go:36", "field_fn_key1": "field_fn_value1"}
2023-02-21 01:11:08.972	error	project	wlog/logger.go:173	{"msg": "Error", "caller": "main.go:37", "field_fn_key1": "field_fn_value1"}
2023-02-21 01:11:08.972	dpanic	project	wlog/logger.go:182	{"msg": "DPanic", "caller": "main.go:38", "field_fn_key1": "field_fn_value1"}
github.com/wyy-go/wlog.(*Log).DPanic
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/logger.go:182
github.com/wyy-go/wlog.DPanic
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/default.go:106
main.main
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/example/main.go:38
runtime.main
	D:/gvm/go/src/runtime/proc.go:250
2023-02-21 01:11:08.973	debug	project	wlog/logger.go:206	{"msg": "Debugf: debug", "caller": "main.go:40", "field_fn_key1": "field_fn_value1"}
2023-02-21 01:11:08.973	info	project	wlog/logger.go:214	{"msg": "Infof: info", "caller": "main.go:41", "field_fn_key1": "field_fn_value1"}
2023-02-21 01:11:08.973	warn	project	wlog/logger.go:222	{"msg": "Warnf: warn", "caller": "main.go:42", "field_fn_key1": "field_fn_value1"}
2023-02-21 01:11:08.973	info	project	wlog/logger.go:214	{"msg": "Infof: info", "caller": "main.go:43", "field_fn_key1": "field_fn_value1"}
2023-02-21 01:11:08.974	error	project	wlog/logger.go:230	{"msg": "Errorf: error", "caller": "main.go:44", "field_fn_key1": "field_fn_value1"}
2023-02-21 01:11:08.974	dpanic	project	wlog/logger.go:239	{"msg": "DPanicf: dPanic", "caller": "main.go:45", "field_fn_key1": "field_fn_value1"}
github.com/wyy-go/wlog.(*Log).DPanicf
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/logger.go:239
github.com/wyy-go/wlog.DPanicf
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/default.go:128
main.main
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/example/main.go:45
runtime.main
	D:/gvm/go/src/runtime/proc.go:250
2023-02-21 01:11:08.974	debug	project	wlog/logger.go:268	{"msg": "Debugw", "caller": "main.go:47", "field_fn_key1": "field_fn_value1", "Debugw": "w", "11111111": "2222222222"}
2023-02-21 01:11:08.974	info	project	wlog/logger.go:277	{"msg": "Infow", "caller": "main.go:48", "field_fn_key1": "field_fn_value1", "Infow": "w"}
2023-02-21 01:11:08.975	warn	project	wlog/logger.go:286	{"msg": "Warnw", "caller": "main.go:49", "field_fn_key1": "field_fn_value1", "Warnw": "w"}
2023-02-21 01:11:08.975	info	project	wlog/logger.go:277	{"msg": "Infow", "caller": "main.go:50", "field_fn_key1": "field_fn_value1", "Infow": "w"}
2023-02-21 01:11:08.975	error	project	wlog/logger.go:295	{"msg": "Errorw", "caller": "main.go:51", "field_fn_key1": "field_fn_value1", "Errorw": "w"}
2023-02-21 01:11:08.975	dpanic	project	wlog/logger.go:305	{"msg": "DPanicw", "caller": "main.go:52", "field_fn_key1": "field_fn_value1", "DPanicw": "w"}
github.com/wyy-go/wlog.(*Log).DPanicw
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/logger.go:305
github.com/wyy-go/wlog.DPanicw
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/default.go:159
main.main
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/example/main.go:52
runtime.main
	D:/gvm/go/src/runtime/proc.go:250
2023-02-21 01:11:08.975	panic	project	wlog/logger.go:190	{"msg": "Panic", "caller": "main.go:55", "field_fn_key1": "field_fn_value1"}
github.com/wyy-go/wlog.(*Log).Panic
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/logger.go:190
github.com/wyy-go/wlog.Panic
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/default.go:109
main.main.func1
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/example/main.go:55
main.shouPanic
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/example/main.go:89
main.main
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/example/main.go:54
runtime.main
	D:/gvm/go/src/runtime/proc.go:250
2023-02-21 01:11:08.976	panic	project	wlog/logger.go:247	{"msg": "Panicf: panic", "caller": "main.go:58", "field_fn_key1": "field_fn_value1"}
github.com/wyy-go/wlog.(*Log).Panicf
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/logger.go:247
github.com/wyy-go/wlog.Panicf
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/default.go:131
main.main.func2
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/example/main.go:58
main.shouPanic
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/example/main.go:89
main.main
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/example/main.go:57
runtime.main
	D:/gvm/go/src/runtime/proc.go:250
2023-02-21 01:11:08.976	panic	project	wlog/logger.go:314	{"msg": "Panicw: %s", "caller": "main.go:61", "field_fn_key1": "field_fn_value1", "panic": "w"}
github.com/wyy-go/wlog.(*Log).Panicw
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/logger.go:314
github.com/wyy-go/wlog.Panicw
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/default.go:163
main.main.func3
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/example/main.go:61
main.shouPanic
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/example/main.go:89
main.main
	C:/Users/wangyangyang/Desktop/qwqqq/wlog/example/main.go:60
runtime.main
	D:/gvm/go/src/runtime/proc.go:250
2023-02-21 01:11:08.977	debug	project	wlog/logger.go:149	{"aa": "bb", "msg": "debug with", "caller": "main.go:64", "field_fn_key1": "field_fn_value1"}
2023-02-21 01:11:08.977	debug	project.another	wlog/logger.go:149	{"msg": "debug named", "caller": "main.go:66", "field_fn_key1": "field_fn_value1"}
2023-02-21 01:11:08.977	debug	project	wlog/logger.go:149	{"msg": "with context", "caller": "main.go:70", "field_fn_key1": "field_fn_value1", "field_fn_key2": "field_fn_value2"}
2023-02-21 01:11:08.977	debug	project	wlog/logger.go:149	{"msg": "with field fn", "caller": "main.go:74", "field_fn_key1": "field_fn_value1", "field_fn_key3": "field_fn_value3"}
2023-02-21 01:11:08.977	debug	project	wlog/logger.go:149	{"aaaa": {"xx": "yy", "bbbbbb": {"dd": "gg", "msg": "", "caller": "main.go:76", "field_fn_key1": "field_fn_value1"}}}

```