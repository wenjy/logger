package logger

import (
	"fmt"
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	zap         *zap.Logger
	atomicLevel zap.AtomicLevel
}

//日志级别映射
var zapLevelMap = map[Level]zapcore.Level{
	DEBUG: zap.DebugLevel,
	INFO:  zap.InfoLevel,
	WARN:  zap.WarnLevel,
	ERROR: zap.ErrorLevel,
	FATAL: zap.FatalLevel,
}

//注册zap日志库
func RegisterZapLogger(level Level, hook *lumberjack.Logger, debug bool) *zapLogger {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "file",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 短路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	var writes = []zapcore.WriteSyncer{zapcore.AddSync(hook)}
	var options []zap.Option
	options = append(options, zap.AddStacktrace(zap.ErrorLevel))

	//开启debug模式支持输出到控制台，以及记录函数调用信息
	if debug {
		writes = append(writes, zapcore.AddSync(os.Stdout))
		options = append(options, zap.AddCaller())
		options = append(options, zap.AddCallerSkip(2))
	}

	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zapLevelMap[level])

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(writes...),
		atomicLevel,
	)

	// 构造日志
	return &zapLogger{zap: zap.New(core, options...), atomicLevel: atomicLevel}
}

// 设置日志级别
func (zapLog *zapLogger) SetLevel(level Level) {
	zapLog.atomicLevel.SetLevel(zapLevelMap[level])
}

func (zapLog *zapLogger) Debugf(msg string, args []Extra) {
	zapLog.zap.Debug(msg, zapLog.transformField(args)...)
}

func (zapLog *zapLogger) Infof(msg string, args []Extra) {
	zapLog.zap.Info(msg, zapLog.transformField(args)...)
}

func (zapLog *zapLogger) Warnf(msg string, args []Extra) {
	zapLog.zap.Warn(msg, zapLog.transformField(args)...)
}

func (zapLog *zapLogger) Errorf(msg string, args []Extra) {
	zapLog.zap.Error(msg, zapLog.transformField(args)...)
}

func (zapLog *zapLogger) Fatalf(msg string, args []Extra) {
	zapLog.zap.Fatal(msg, zapLog.transformField(args)...)
}

//转换扩展信息对象为zap的Field列表
func (zapLog *zapLogger) transformField(args []Extra) (result []zap.Field) {
	result = make([]zap.Field, len(args))
	for index, arg := range args {
		switch arg.Value.(type) {
		case bool:
			result[index] = zap.Bool(arg.Key, arg.Value.(bool))
		case string:
			result[index] = zap.String(arg.Key, arg.Value.(string))
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			result[index] = zap.String(arg.Key, fmt.Sprintf("%d", arg.Value))
		case []byte:
			result[index] = zap.Binary(arg.Key, arg.Value.([]byte))
		default:
			result[index] = zap.Reflect(arg.Key, arg.Value)
		}
	}
	return result
}
