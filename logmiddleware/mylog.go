package logmiddleware

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var (
	Logger   = logrus.New()
	LogEntry *logrus.Entry
)

func CheckFile(path string, isdir bool) {
	_, err := os.Stat(path)
	if err == nil {
		//文件存在
		return
	}
	if os.IsNotExist(err) {
		//文件不存在
		if isdir {
			err := os.Mkdir(path, 0755)
			if err != nil {
				panic(err)
			}
			return
		} else {
			file, err := os.OpenFile(path, os.O_CREATE, 0644)
			defer file.Close()
			if err != nil {
				panic(err)
			}
			return
		}

	}
	//其他错误
	panic(err)
}

func init() {
	logdir := "logs"
	logpath := logdir + "/log"
	loglink := logdir + "/last.log"

	CheckFile(logdir, true)
	CheckFile(logpath, false)

	src, err := os.OpenFile(logpath, os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	Logger.Out = src
	Logger.SetLevel(logrus.DebugLevel)
	logWriter, _ := rotatelogs.New(
		logpath+"%Y%m%d.log",
		rotatelogs.WithMaxAge(31*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
		rotatelogs.WithLinkName(loglink),
	)
	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter, // info级别使用logWriter写日志
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}
	Hook := lfshook.NewHook(writeMap, &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	Logger.AddHook(Hook)
	LogEntry = logrus.NewEntry(Logger).WithField("service", "test")
}

func MyLog() gin.HandlerFunc {
	//返回handler
	return func(ctx *gin.Context) {
		startTime := time.Now()
		if ctx.Request.UserAgent() == "" {
			ctx.String(403, `
			<html>
			<head><title>403 Forbidden</title></head>
			<body>
			<center><h1>403 Forbidden</h1></center>
			<hr><center>nginx/1.22.0</center>
			</body>
			</html>
			`)
			return
		}
		ctx.Next()
		stopTime := time.Since(startTime)
		spendTime := fmt.Sprintf("%d us", int(math.Ceil(float64(stopTime.Nanoseconds()/1000))))
		statusCode := ctx.Writer.Status()
		remoteIp := ctx.RemoteIP()
		method := ctx.Request.Method
		path := ctx.Request.URL
		useragent := ctx.Request.UserAgent()
		Log := Logger.WithFields(logrus.Fields{
			"Method":     method,
			"URI":        path,
			"StatusCode": statusCode,
			"RemoteIP":   remoteIp,
			"SpendTime":  spendTime,
			"User-Agent": useragent,
			"Host":       ctx.Request.Host,
		})
		if len(ctx.Errors) > 0 { // 矿建内部错误
			Log.Error(ctx.Errors.ByType(gin.ErrorTypePrivate))
		}
		if statusCode >= 500 {
			Log.Error()
		} else if statusCode >= 400 {
			Log.Warn()
		} else {
			Log.Info()
		}
	}
}
