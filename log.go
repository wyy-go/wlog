package wlog

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// New constructs a new Log
func New(opts ...Option) (*zap.Logger, zap.AtomicLevel) {
	c := &Config{}
	for _, opt := range opts {
		opt(c)
	}
	var options []zap.Option

	if c.AddCaller {
		// 添加显示文件名和行号,跳过封装调用层,
		options = append(options, zap.AddCaller(), zap.AddCallerSkip(c.CallerSkip))
	}
	if c.Stack {
		// 栈调用,及使能等级
		options = append(options, zap.AddStacktrace(zap.NewAtomicLevelAt(zap.DPanicLevel))) // 只显示栈的错误等级
	}

	level, err := zap.ParseAtomicLevel(c.Level)
	if err != nil {
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	// 初始化core
	core := zapcore.NewCore(
		toEncoder(c, level), // 设置encoder
		toWriter(c),         // 设置输出
		level,               // 设置日志输出等级
	)
	return zap.New(core, options...), level
}

func toEncoder(c *Config, level zap.AtomicLevel) zapcore.Encoder {
	timeLoyout := "2006-01-02 15:04:05.000"

	encoderConfig := &zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "project",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     zapcore.TimeEncoderOfLayout(timeLoyout),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   customCaller, //zapcore.ShortCallerEncoder,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		//MessageKey:     "msg",
	}
	//if level.Level() == zap.DebugLevel {
	//	encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	//}

	if c.Format == "console" {
		return zapcore.NewConsoleEncoder(*encoderConfig)
	}
	return zapcore.NewJSONEncoder(*encoderConfig)
}

func customCaller(_ zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	var shortPath string
	file, line := caller(7)

	idx := strings.LastIndexByte(file, '/')
	if idx == -1 {
		shortPath = file
	}

	if idx > 0 {
		// Find the penultimate separator.
		idx = strings.LastIndexByte(file[:idx], '/')
		if idx == -1 {
			shortPath = file
		}
	}

	if idx > 0 {
		shortPath = file[idx+1:]
	}

	shortPath = shortPath + ":" + strconv.Itoa(line)
	enc.AppendString(shortPath)
}

func caller(depth int) (file string, line int) {
	d := depth
	_, file, line, _ = runtime.Caller(d)
	if strings.LastIndex(file, "/logger.go") > 0 {
		d++
		_, file, line, _ = runtime.Caller(d)
	}
	if strings.LastIndex(file, "/default.go") > 0 {
		d++
		_, file, line, _ = runtime.Caller(d)
	}
	return file, line
}

func toEncodeLevel(l string) zapcore.LevelEncoder {
	switch l {
	case "LowercaseColorLevelEncoder": // 小写编码器带颜色
		return zapcore.LowercaseColorLevelEncoder
	case "CapitalLevelEncoder": // 大写编码器
		return zapcore.CapitalLevelEncoder
	case "CapitalColorLevelEncoder": // 大写编码器带颜色
		return zapcore.CapitalColorLevelEncoder
	case "LowercaseLevelEncoder": // 小写编码器(默认)
		return zapcore.LowercaseLevelEncoder
	default:
		return zapcore.LowercaseLevelEncoder
	}
}

func toWriter(c *Config) zapcore.WriteSyncer {
	fileWriter := func() zapcore.WriteSyncer {
		return zapcore.AddSync(&lumberjack.Logger{ // 文件切割
			Filename:   filepath.Join(c.Path, c.Filename),
			MaxSize:    c.MaxSize,
			MaxAge:     c.MaxAge,
			MaxBackups: c.MaxBackups,
			LocalTime:  c.LocalTime,
			Compress:   c.Compress,
		})
	}
	stdoutWriter := func() zapcore.WriteSyncer {
		return zapcore.AddSync(os.Stdout)
	}
	customWriter := func(w ...zapcore.WriteSyncer) []zapcore.WriteSyncer {
		ws := make([]zapcore.WriteSyncer, 0, len(c.Writer)+len(w))

		for _, writer := range c.Writer {
			ws = append(ws, zapcore.AddSync(writer))
		}
		for _, writer := range w {
			ws = append(ws, zapcore.AddSync(writer))
		}
		return ws
	}
	switch strings.ToLower(c.Adapter) {
	case "file":
		return fileWriter()
	case "multi":
		return zapcore.NewMultiWriteSyncer(stdoutWriter(), fileWriter())
	case "file-custom":
		return zapcore.NewMultiWriteSyncer(customWriter(fileWriter())...)
	case "console-custom":
		return zapcore.NewMultiWriteSyncer(customWriter(stdoutWriter())...)
	case "multi-custom":
		return zapcore.NewMultiWriteSyncer(customWriter(stdoutWriter(), fileWriter())...)
	case "custom":
		ws := customWriter()
		if len(ws) == 0 {
			return stdoutWriter()
		}
		if len(ws) == 1 {
			return ws[0]
		}
		return zapcore.NewMultiWriteSyncer(ws...)
	default: // console
		return stdoutWriter()
	}
}
