package logger

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/Chayakorn2002/pms-classroom-backend/domain/exceptions"
)

var Slog *slog.Logger
var m sync.Mutex

type LogConfig struct {
	Env             string `mapstructure:"env"` // Ask P'Mind: Do we need to separate the log config for each environment?
	Level           string `mapstructure:"level"`
	UseJsonEncoder  bool   `mapstructure:"useJsonEncoder"`
	StacktraceLevel string `mapstructure:"stacktraceLevel"`
	FileEnabled     bool   `mapstructure:"fileEnabled"`
	FileSize        int    `mapstructure:"fileSize"`
	FilePath        string `mapstructure:"filePath"`
	FileCompress    bool   `mapstructure:"fileCompress"`
	MaxAge          int    `mapstructure:"maxAge"`
	MaxBackups      int    `mapstructure:"maxBackups"`
}

func InitLogger() {
	m.Lock()
	defer m.Unlock()

	Slog = newZapLogger()

	slog.SetDefault(Slog)
	CompileCanonicalLogTemplate()
	slog.InfoContext(context.Background(), "Logger initialized")
}

var canonicalLogTemplate *template.Template

type Level int

const (
	Debug Level = 1 << iota
	Info
	Warn
	Error
)

type CanonicalLog struct {
	Transport string
	Traffic   string
	Method    string
	Status    int
	Path      string
	Duration  time.Duration
	Message   string
}

func CompileCanonicalLogTemplate() {
	logTemplate := "[{{.Transport}}][{{.Traffic}}] {{.Method}} {{.Status}} {{.Path}} {{.Duration}} - {{.Message}}"
	compiled, err := template.New("log_template").Parse(logTemplate)
	if err != nil {
		panic(err)
	}
	canonicalLogTemplate = compiled
}

func GetCanonicalLogTemplate() (*template.Template, error) {
	if canonicalLogTemplate != nil {
		return canonicalLogTemplate, nil
	}
	return nil, errors.New("canonicalLogTemplate is nil")
}

func CanonicalLogger(ctx context.Context, slogger slog.Logger, level Level, request []byte, response []byte, err error, cannonicalLog CanonicalLog, metadata []any) {
	// log the cannonical log

	var fields []any
	// append request log
	var jsonObj map[string]interface{}
	if unmarshalErr := json.Unmarshal(request, &jsonObj); unmarshalErr != nil {
		fields = append(fields, slog.String("request", string(request)))
	} else {
		fields = append(fields, slog.Any("request", jsonObj))
	}

	// append response log
	if err != nil {
		level = Error
		cErr, ok := err.(*exceptions.ExceptionError)
		if ok && cErr != nil {
			if cErr.StackErrors != nil {
				stackTrace := exceptions.GetStackField(cErr.StackErrors)
				stackTraceParts := strings.Split(stackTrace.Stack, "\n\t")
				if len(stackTraceParts) > 6 {
					stackTrace.Stack = strings.Join(stackTraceParts[:6], "\n\t")
				}
				fields = append(fields, slog.Group("error",
					slog.String("kind", stackTrace.Kind),
					slog.String("message", stackTrace.Message),
					slog.String("stack", stackTrace.Stack),
				))

			}
			fields = append(fields, slog.Group("response",
				slog.Int("status_code", cErr.APIStatusCode),
				slog.Any("data", nil),
				slog.Group("error",
					slog.Int("code", int(cErr.Code)),
					slog.String("message", cErr.GlobalMessage),
					slog.String("debug_message", cErr.DebugMessage),
				)))
			cannonicalLog.Message = cErr.DebugMessage
		} else {
			// This is the case when the error is not an instance of ExceptionError
			var jsonObj map[string]interface{}
			if err := json.Unmarshal(response, &jsonObj); err != nil {
				fields = append(fields, slog.Group("response",
					slog.Int("status_code", cannonicalLog.Status),
					slog.String("data", string(response)),
				))
			} else {
				fields = append(fields, slog.Group("response",
					slog.Int("status_code", cannonicalLog.Status),
					slog.Any("data", jsonObj),
				))
			}
		}
	} else {
		level = Info
		var jsonObj map[string]interface{}
		if err := json.Unmarshal(response, &jsonObj); err != nil {
			fields = append(fields, slog.Group("response",
				slog.Int("status_code", 1000),
				slog.String("data", string(response)),
				slog.Any("error", nil),
			))
		} else {
			fields = append(fields, slog.Group("response",
				slog.Int("status_code", 1000),
				slog.Any("data", jsonObj),
				slog.Any("error", nil),
			))
		}
	}

	// append md log
	fields = append(fields,
		slog.String("logger_name", "canonical"),
		slog.Group("md", metadata...),
	)

	var logMsgBuilder strings.Builder
	var logMsg string
	logTmpl, logTmplErr := GetCanonicalLogTemplate()
	if logTmplErr != nil {
		logMsg = "failed to get cannonical log template"
	} else {
		executeErr := logTmpl.Execute(&logMsgBuilder, cannonicalLog)
		if executeErr != nil {
			logMsg = "failed to execute cannonical log template"
		} else {
			logMsg = logMsgBuilder.String()
		}
	}
	switch level {
	case Debug:
		slogger.DebugContext(ctx, logMsg, fields...)
	case Info:
		slogger.InfoContext(ctx, logMsg, fields...)
	case Warn:
		slogger.WarnContext(ctx, logMsg, fields...)
	case Error:
		slogger.ErrorContext(ctx, logMsg, fields...)
	default:
		slogger.ErrorContext(ctx, logMsg, fields...)
	}
}
