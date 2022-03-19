package logger

import (
	"fmt"

	"github.com/natefinch/lumberjack"
)

func NewLogger(filename string, maxSize, maxDay int, driver, level string, debug bool) Logger {
	hook := lumberjack.Logger{
		Filename:   filename, // 日志文件路径
		MaxSize:    maxSize,  // 每个日志文件保存的大小 单位:M
		MaxAge:     maxDay,   // 文件最多保存多少天
		MaxBackups: 3,        // 日志文件最多保存多少个备份
		Compress:   false,    // 是否压缩
		LocalTime:  true,
	}

	var logger Logger
	switch driver {
	case "zap":
		if level, ok := LevelMap[level]; !ok {
			panic(fmt.Errorf("Unknown log level.\n"))
		} else {
			logger = RegisterZapLogger(level, &hook, debug)
		}
	default:
		panic(fmt.Errorf("Unknown log driver.\n"))
	}
	return logger
}
