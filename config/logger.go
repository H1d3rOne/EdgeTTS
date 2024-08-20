package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

// 设置日志级别
var LogConf = struct {
	Dir     string `yaml:"dir"`
	Name    string `yaml:"name"`
	Level   string `yaml:"level"`
	MaxSize int    `yaml:"max_size"`
}{
	Dir:     "./logs",
	Name:    "yourlogname.log",
	Level:   "trace",
	MaxSize: 100,
}

// 初始化logger
func InitLogger() error {
	//创建一个logger对象
	var log = logrus.New()
	//设置日志样式
	log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
	})
	switch LogConf.Level {
	case "trace":
		log.SetLevel(logrus.TraceLevel)
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	case "fatal":
		log.SetLevel(logrus.FatalLevel)
	case "panic":
		log.SetLevel(logrus.PanicLevel)
	}
	log.SetReportCaller(true) // 打印文件、行号和主调函数。

	// 实现日志滚动。
	logger := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%v/%v", LogConf.Dir, LogConf.Name), // 日志输出文件路径。
		MaxSize:    LogConf.MaxSize,                                 // 日志文件最大 size(MB)，缺省 100MB。
		MaxBackups: 10,                                              // 最大过期日志保留的个数。
		MaxAge:     30,                                              // 保留过期文件的最大时间间隔，单位是天。
		LocalTime:  true,                                            // 是否使用本地时间来命名备份的日志。
	}

	// 同时输出到标准输出与文件。
	log.SetOutput(io.MultiWriter(logger, os.Stdout))
	return nil

}
