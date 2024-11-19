package logger

import (
	"os"

	"log/slog"

	// "github.com/Orbit-Digital-Company-Limited-MP/SH002-shared-library/utils/runtime"
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func newZapLogger() *slog.Logger {
	// STEP 0: Get the log profile based on env
	log := LogConfig{
		Env:             "local",
		Level:           "debug",
		UseJsonEncoder:  false,
		StacktraceLevel: "error",
		FileEnabled:     false,
		FilePath:        "logs/app.log",
		FileSize:        100, // megabytes
		FileCompress:    true,
		MaxAge:          30, // days
		MaxBackups:      3,  // number of log files
	}

	// STEP 1: Get the log level
	zapLogLevel := getZapLogLevel(log.Level)
	//stacktraceLogLevel := getZapLogLevel(log.StacktraceLevel)

	// STEP 2: Set up the file writer
	lumberjackLogger := &lumberjack.Logger{
		Filename:   log.FilePath,
		MaxSize:    log.FileSize,     // megabytes
		MaxBackups: log.MaxBackups,   // number of log files
		MaxAge:     log.MaxAge,       // days
		Compress:   log.FileCompress, // disabled by default
	}

	fileWriter := zapcore.AddSync(lumberjackLogger)

	// STEP 3: Set up the encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// STEP 4: Set up the encoder for the file before changing it for the console
	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// STEP 5: Change the time format for the console
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	// STEP 6: Set up the core

	zapCoreList := []zapcore.Core{}
	if log.FileEnabled {
		zapCoreList = append(zapCoreList, zapcore.NewCore(jsonEncoder, fileWriter, zapLogLevel))
	}

	if log.UseJsonEncoder {
		zapCoreList = append(zapCoreList, zapcore.NewCore(jsonEncoder, zapcore.AddSync(os.Stdout), zapLogLevel))
	}

	var core zapcore.Core
	// Set up the console for default
	if len(zapCoreList) == 0 {
		core = zapcore.NewTee(zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapLogLevel))
	} else {
		// Set up the console for the rest
		core = zapcore.NewTee(zapCoreList...)
	}

	// STEP 7: Set up the logger
	//logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(stacktraceLogLevel), zap.AddCallerSkip(skip))

	// STEP 7: Set up the slog logger
	logger := slog.New(zapslog.NewHandler(core, &zapslog.HandlerOptions{
		AddSource: true,
	}))
	return logger
}

func getZapLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}
