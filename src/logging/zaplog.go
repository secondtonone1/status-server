package logging

import (
	"fmt"
	"path/filepath"
	"status-server/config"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	logger      *zap.Logger
	sugarLogger *zap.SugaredLogger
	level       zapcore.Level
}

func NewZapLogger(conf *config.LogConf) *ZapLogger {

	filename := ""
	ext := conf.Extname
	if ext == "" {
		filename = filepath.Join(conf.LogDir, fmt.Sprintf("%s-%s.log", conf.Project, conf.Name))
	} else {
		//ext = strconv.FormatInt(rand.Int63(), 10)
		filename = filepath.Join(conf.LogDir, fmt.Sprintf("%s-%s.log.%s", conf.Project, conf.Name, ext))
	}

	maxsize := conf.MaxSize
	if maxsize <= 0 {
		maxsize = 1024
	}

	maxnum := conf.MaxNum
	if maxnum <= 0 {
		maxnum = 10
	}

	//hook := zapwrapper.Logger{
	//	Filename:   filename, // 日志文件路径
	//	MaxSize:    maxsize,  // megabytes
	//	MaxBackups: maxnum,   // 最多保留5个备份
	//	MaxAge:     30,       //days
	//	Compress:   true,     // 是否压缩 disabled by default
	//	LocalTime:  true,
	//}

	if conf.MaxDay <= 0 {
		conf.MaxDay = 7
	}

	if conf.RotateSeconds <= 10 {
		conf.RotateSeconds = 3600
	}

	hook, err := rotatelogs.New(
		filename+".%Y%m%d%H%M%S",
		//rotatelogs.WithLinkName(baseLogPaht),      						// 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(time.Hour*24*time.Duration(conf.MaxDay)),             // 文件最大保存时间
		rotatelogs.WithRotationTime(time.Second*time.Duration(conf.RotateSeconds)), // 日志切割时间间隔
	)

	if err != nil {
		panic("new rotatelogs fail")
	}

	w := zapcore.AddSync(hook)

	var level zapcore.Level
	switch conf.LogLevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "error":
		level = zap.ErrorLevel
	case "dpanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}
	encoderConfig := zap.NewProductionEncoderConfig()
	if conf.NotPrintLogTime {
		encoderConfig.LevelKey = ""
		encoderConfig.TimeKey = ""
	} else {
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		w,
		level,
	)

	logger := zap.New(core)
	//logger.Info("DefaultLogger init success")

	return &ZapLogger{
		logger:      logger,
		sugarLogger: logger.Sugar(),
		level:       level,
	}
}

func (zl *ZapLogger) Info(args ...interface{}) {
	if zl.level <= zap.InfoLevel {
		zl.sugarLogger.Info(args...)
	}
}

func (zl *ZapLogger) Error(args ...interface{}) {
	if zl.level <= zap.ErrorLevel {
		zl.sugarLogger.Error(args...)
	}
}

func (zl *ZapLogger) Debug(args ...interface{}) {
	if zl.level <= zap.DebugLevel {
		zl.sugarLogger.Debug(args...)
	}
}

func (zl *ZapLogger) Warn(args ...interface{}) {
	if zl.level <= WarnLevel {
		zl.sugarLogger.Warn(args...)
	}
}

func (zl *ZapLogger) DPanic(args ...interface{}) {
	if zl.level <= DPanicLevel {
		zl.sugarLogger.DPanic(args...)
	}
}

func (zl *ZapLogger) Panic(args ...interface{}) {
	if zl.level <= PanicLevel {
		zl.sugarLogger.Panic(args...)
	}
}

func (zl *ZapLogger) Fatal(args ...interface{}) {
	if zl.level <= FatalLevel {
		zl.sugarLogger.Fatal(args...)
	}
}

func (zl *ZapLogger) Infof(format string, args ...interface{}) {
	if zl.level <= zap.InfoLevel {
		zl.sugarLogger.Infof(format, args...)
	}
}

func (zl *ZapLogger) Errorf(format string, args ...interface{}) {
	if zl.level <= zap.ErrorLevel {
		zl.sugarLogger.Errorf(format, args...)
	}
}

func (zl *ZapLogger) Debugf(format string, args ...interface{}) {
	if zl.level <= zap.DebugLevel {
		zl.sugarLogger.Debugf(format, args...)
	}
}

func (zl *ZapLogger) Warnf(format string, args ...interface{}) {
	if zl.level <= WarnLevel {
		zl.sugarLogger.Warnf(format, args...)
	}
}

func (zl *ZapLogger) DPanicf(format string, args ...interface{}) {
	if zl.level <= DPanicLevel {
		zl.sugarLogger.DPanicf(format, args...)
	}
}

func (zl *ZapLogger) Panicf(format string, args ...interface{}) {
	if zl.level <= PanicLevel {
		zl.sugarLogger.Panicf(format, args...)
	}
}

func (zl *ZapLogger) Fatalf(format string, args ...interface{}) {
	if zl.level <= FatalLevel {
		zl.sugarLogger.Fatalf(format, args...)
	}
}

func (zl *ZapLogger) SetLevel(l string) {
	switch l {
	case "debug":
		zl.level = zap.DebugLevel
	case "info":
		zl.level = zap.InfoLevel
	case "warn":
		zl.level = zap.WarnLevel
	case "error":
		zl.level = zap.ErrorLevel
	case "dpanic":
		zl.level = zap.DPanicLevel
	case "panic":
		zl.level = zap.PanicLevel
	case "fatal":
		zl.level = zap.FatalLevel
	default:
		zl.level = zap.InfoLevel
	}
}
