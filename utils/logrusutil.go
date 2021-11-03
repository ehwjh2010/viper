package utils

import (
	"ginLearn/client/setting"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"time"
)

var Log *logrus.Logger
var f *os.File

//InitLogrus logrus初始化设置
func InitLogrus(application string, logConfig *setting.LogConfig) error {
	var writers []io.Writer

	if IsNotEmptyStr(logConfig.LogPath) {
		//确保日志目录存在
		dirLogPath := PathJoin(logConfig.LogPath, application)
		err := MakeDirs(dirLogPath)
		if err != nil {
			//log.Fatalf("Access log dir failed! err: %v", err)
			return err
		}

		//确保日志文件存在
		logFilePath := PathJoin(dirLogPath, "application.log")
		f, err = OpenFileWithAppend(logFilePath)
		if err != nil {
			//log.Fatalf("Access log file failed! err: %v", err)
			return err
		}
	}

	if f != nil {
		log.Println("Log use file writer")
		writers = append(writers, f)
	}

	if logConfig.EnableLogConsole {
		log.Println("Log use console writer")
		writers = append(writers, os.Stdout)
	}

	//实例化
	logger := logrus.New()

	//设置输出
	if len(writers) == 0 {
		log.Println("No set log writer, User console as default writer!!!")
		logger.SetOutput(os.Stdout)
	} else {
		logger.SetOutput(io.MultiWriter(writers...))
	}

	//设置gin框架相关日志输出
	gin.DefaultWriter = logger.Out
	gin.DefaultErrorWriter = logger.Out

	//设置日志级别
	level, err := logrus.ParseLevel(logConfig.Level)
	if err != nil {
		//logger.Fatalf("logger level convert failed!, err: %v", err)
		return err
	}
	logger.SetLevel(level)

	//设置日志格式
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	//添加打印日志所在文件以及行数, 比较影响性能, 是否使用自行决定
	if logConfig.AccessMethodRow {
		logger.SetReportCaller(true)
	}

	Log = logger

	return nil
}

func CloseLogFile() error {

	if f == nil {
		return nil
	}

	err := f.Close()

	if err == nil {
		log.Println("Close log file success")
	} else {
		log.Println("Close log file failed")
	}

	return err
}