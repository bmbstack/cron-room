package job

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	. "github.com/onestack/cron-room/helper"
	"gopkg.in/natefinch/lumberjack.v2"
	"time"
	"fmt"
	"runtime"
)

type BaseJob struct {
	Logger *zap.Logger
}

func CreateJob(name string) *BaseJob {
	filename := fmt.Sprintf("%s%s.log", GetEnv().LogFilePathPrefix, name)
	ws := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filename,
		MaxSize:    500, // megabytes
		MaxBackups: 3,   // backup
		MaxAge:     30,  // days
	})
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = ""
	coreInstance := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		ws,
		zap.InfoLevel,
	)
	return &BaseJob{Logger: zap.New(coreInstance)}
}

func (this BaseJob) WriteInfo(msg string, data string) {
	this.Logger.Info(msg,
		zap.String(LogKeySource, "[协程]"),
		zap.String(LogKeyTime, time.Now().Format(DateFullLayout)),
		zap.String(LogKeyData, data),
		zap.Int(LogKeyGoroutineNum, runtime.NumGoroutine()),
	)
}



