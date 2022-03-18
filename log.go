package logger

import (
	"fmt"

	"github.com/natefinch/lumberjack"
)

type Log struct {
	logger Logger
}

func NewLogger(filename string, maxSize, maxDay int, driver, level string, debug bool) *Log {
	hook := lumberjack.Logger{
		Filename:   filename, // 日志文件路径
		MaxSize:    maxSize,  // 每个日志文件保存的大小 单位:M
		MaxAge:     maxDay,   // 文件最多保存多少天
		MaxBackups: 3,        // 日志文件最多保存多少个备份
		Compress:   false,    // 是否压缩
		LocalTime:  true,
	}

	logger := &Log{}
	switch driver {
	case "zap":
		if level, ok := LevelMap[level]; !ok {
			panic(fmt.Errorf("Unknown log level.\n"))
		} else {
			logger.logger = RegisterZapLogger(level, &hook, debug)
		}
	default:
		panic(fmt.Errorf("Unknown log driver.\n"))
	}
	return logger
}

func (l *Log) Debugf(msg string, args ...Extra) {
	l.logger.Debugf(msg, args)
}

func (l *Log) Infof(msg string, args ...Extra) {
	l.logger.Infof(msg, args)
}

func (l *Log) Warnf(msg string, args ...Extra) {
	l.logger.Warnf(msg, args)
}

func (l *Log) Errorf(msg string, args ...Extra) {
	l.logger.Errorf(msg, args)
}

func (l *Log) Fatalf(msg string, args ...Extra) {
	l.logger.Fatalf(msg, args)
}
