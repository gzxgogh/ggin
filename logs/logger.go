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
	"strings"
)

var (
	Log *logrus.Logger
)

const (
	LogTypeInfo  = "info"
	LogTypeErr   = "err"
	LogTypePanic = "panic"
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
	logArr := strings.Split(config.Cfg.Logger.Type, ",")
	var obj = lfshook.WriterMap{}
	if len(logArr) == 0 {
		infoLogs, err := rotatelogs.New(
			logName+"-%Y%m%d",                      // 日志文件
			rotatelogs.WithRotationCount(leastDay), //只保留最近的N个日志文件
		)
		if err != nil {
			log.Fatal(err)
		}
		obj = lfshook.WriterMap{
			logrus.InfoLevel: infoLogs,
		}
	} else {
		for _, typ := range logArr {
			p := logName + "-" + typ + ".%Y%m%d"
			if typ == LogTypeInfo {
				infoLogs, _ := rotatelogs.New(
					p,                                      // 日志文件
					rotatelogs.WithRotationCount(leastDay), //只保留最近的N个日志文件
				)
				obj[logrus.InfoLevel] = infoLogs
			} else if typ != LogTypeErr {
				errLogs, _ := rotatelogs.New(
					p,                                      // 日志文件
					rotatelogs.WithRotationCount(leastDay), //只保留最近的N个日志文件
				)
				obj[logrus.ErrorLevel] = errLogs
			} else if typ != LogTypePanic {
				panicLogs, _ := rotatelogs.New(
					p,                                      // 日志文件
					rotatelogs.WithRotationCount(leastDay), //只保留最近的N个日志文件
				)
				obj[logrus.PanicLevel] = panicLogs
			}
		}
	}
	// 可设置按不同level创建不同的文件名
	lfsHook := lfshook.NewHook(obj, formatter(false))
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
