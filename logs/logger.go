package logs

import (
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/gzxgogh/ggin/config"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"log"
	"path/filepath"
	"runtime"
)

var (
	Log *logrus.Logger
)

// InitLogger InitLogger
func InitLogger() {
	Log = logrus.New()
	hook := NewLfsHook(filepath.Join(config.Cfg.Logger.Path, config.Cfg.App.Name), 7)
	Log.AddHook(hook)
	Log.SetFormatter(formatter(true))
	Log.SetReportCaller(true)
}

func formatter(isConsole bool) *nested.Formatter {
	fmtter := &nested.Formatter{
		HideKeys:        true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerFirst:     true,
		CustomCallerFormatter: func(frame *runtime.Frame) string {
			funcInfo := runtime.FuncForPC(frame.PC)
			if funcInfo == nil {
				return "error during runtime.FuncForPC"
			}
			fullPath, line := funcInfo.FileLine(frame.PC)
			return fmt.Sprintf(" [%v:%v]", filepath.Base(fullPath), line)
		},
	}
	if runtime.GOOS == "windows" {
		fmtter.NoColors = true
	} else {
		fmtter.NoColors = false
	}
	return fmtter
}

// NewLfsHook NewLfsHook
func NewLfsHook(logName string, leastDay uint) logrus.Hook {
	infoLogs, err := rotatelogs.New(
		logName+"-info.%Y%m%d",                 // 日志文件
		rotatelogs.WithRotationCount(leastDay), //只保留最近的N个日志文件
	)
	if err != nil {
		log.Fatal(err)
	}
	errLogs, err := rotatelogs.New(
		logName+"-err.%Y%m%d",                  // 日志文件
		rotatelogs.WithRotationCount(leastDay), //只保留最近的N个日志文件
	)
	if err != nil {
		log.Fatal(err)
	}
	panicLogs, err := rotatelogs.New(
		logName+"-panic.%Y%m%d",                // 日志文件
		rotatelogs.WithRotationCount(leastDay), //只保留最近的N个日志文件
	)
	if err != nil {
		log.Fatal(err)
	}

	// 可设置按不同level创建不同的文件名
	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.InfoLevel:  infoLogs,
		logrus.ErrorLevel: errLogs,
		logrus.PanicLevel: panicLogs,
	}, formatter(false))
	return lfsHook
}

func Error(args ...interface{}) {
	Log.Log(logrus.ErrorLevel, args...)
}

func Info(args ...interface{}) {
	Log.Log(logrus.InfoLevel, args...)
}

func Panic(args ...interface{}) {
	Log.Log(logrus.PanicLevel, args...)
}
