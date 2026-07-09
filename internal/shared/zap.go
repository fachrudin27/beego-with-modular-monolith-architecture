package shared

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger = zap.NewNop()

const maxLogBodySize = 4096

var sensitiveLogFields = map[string]struct{}{
	"access_token":  {},
	"authorization": {},
	"password":      {},
	"refresh_token": {},
	"token":         {},
}

type zapLoggerConfig struct {
	Development      bool     `json:"development"`
	Encoding         string   `json:"encoding"`
	Level            string   `json:"level"`
	OutputPaths      []string `json:"outputPaths"`
	ErrorOutputPaths []string `json:"errorOutputPaths"`
}

func InitZapLogger(config string) error {
	cfg, err := buildZapConfig(config)
	if err != nil {
		return err
	}

	zapLogger, err := cfg.Build()
	if err != nil {
		return err
	}

	logger = zapLogger
	return nil
}

func SyncZapLogger() {
	_ = logger.Sync()
}

func Logger() *zap.Logger {
	return logger
}

func buildZapConfig(config string) (zap.Config, error) {
	opts := zapLoggerConfig{
		Development:      false,
		Encoding:         "json",
		Level:            "info",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	if config != "" && config != "{}" {
		if err := json.Unmarshal([]byte(config), &opts); err != nil {
			return zap.Config{}, err
		}
	}

	level := zap.NewAtomicLevel()
	if err := level.UnmarshalText([]byte(opts.Level)); err != nil {
		return zap.Config{}, fmt.Errorf("invalid zap level %q: %w", opts.Level, err)
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeDuration = zapcore.StringDurationEncoder

	if opts.Development {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	if len(opts.OutputPaths) == 0 {
		opts.OutputPaths = []string{"stdout"}
	}
	if len(opts.ErrorOutputPaths) == 0 {
		opts.ErrorOutputPaths = []string{"stderr"}
	}

	return zap.Config{
		Level:             level,
		Development:       opts.Development,
		DisableCaller:     true,
		DisableStacktrace: !opts.Development,
		Sampling:          nil,
		Encoding:          opts.Encoding,
		EncoderConfig:     encoderConfig,
		OutputPaths:       opts.OutputPaths,
		ErrorOutputPaths:  opts.ErrorOutputPaths,
	}, nil
}

// error, warn, warning, debug
func ZapLogger(status, title, service, requestID, url string, requestBody []byte, responseBody []byte) {

	_, file, line, ok := runtime.Caller(1)

	position := "unknown:0"
	if ok {
		position = formatRelativePath(file, line)
	}

	fields := []zap.Field{
		zap.String("service", service),
		zap.String("position", position),
		zap.String("request_id", requestID),
		zap.String("url", url),
		zap.String("request", sanitizeLogBody(requestBody)),
		zap.String("response", sanitizeLogBody(responseBody)),
	}

	switch strings.ToLower(status) {
	case "error":
		logger.Error(title, fields...)
	case "warn", "warning":
		logger.Warn(title, fields...)
	case "debug":
		logger.Debug(title, fields...)
	default:
		logger.Info(title, fields...)
	}
}

func sanitizeLogBody(body []byte) string {
	if len(body) == 0 {
		return ""
	}

	var value any
	if err := json.Unmarshal(body, &value); err == nil {
		value = redactLogValue(value)
		redacted, err := json.Marshal(value)
		if err == nil {
			return limitLogBody(redacted)
		}
	}

	return limitLogBody(body)
}

func redactLogValue(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		for key, val := range typed {
			if _, sensitive := sensitiveLogFields[strings.ToLower(key)]; sensitive {
				typed[key] = "[REDACTED]"
				continue
			}
			typed[key] = redactLogValue(val)
		}
		return typed
	case []any:
		for idx, val := range typed {
			typed[idx] = redactLogValue(val)
		}
		return typed
	default:
		return value
	}
}

func limitLogBody(body []byte) string {
	body = bytes.TrimSpace(body)
	if len(body) <= maxLogBodySize {
		return string(body)
	}

	return string(body[:maxLogBodySize]) + "...[TRUNCATED]"
}
