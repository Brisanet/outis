package outis

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogFields representa os campos do log
type LogFields map[string]interface{}

// LogLevel representa o level do log
type LogLevel string

const (
	// DebugLevel representa o nível de debug do Log
	DebugLevel LogLevel = "DebugLevel"
	// InfoLevel representa o nível de info do Log
	InfoLevel LogLevel = "InfoLevel"
	// WarnLevel representa o nível de warn do Log
	WarnLevel LogLevel = "WarnLevel"
	// ErrorLevel representa o nível de error do Log
	ErrorLevel LogLevel = "ErrorLevel"
	// DPanicLevel representa o nível de dpanic do Log
	DPanicLevel LogLevel = "DPanicLevel"
	// PanicLevel representa o nível de panic do Log
	PanicLevel LogLevel = "PanicLevel"
	// FatalLevel representa o nível de fatal do Log
	FatalLevel LogLevel = "FatalLevel"
)

func convertLevel(levelIn LogLevel) (zapcore.Level, bool) {
	zapLevel, hasZapLevel := map[LogLevel]zapcore.Level{
		DebugLevel:  zapcore.DebugLevel,
		InfoLevel:   zapcore.InfoLevel,
		WarnLevel:   zapcore.WarnLevel,
		ErrorLevel:  zapcore.ErrorLevel,
		DPanicLevel: zapcore.DPanicLevel,
		PanicLevel:  zapcore.PanicLevel,
		FatalLevel:  zapcore.FatalLevel,
	}[levelIn]

	return zapLevel, hasZapLevel
}

// LogOptions representa as opções de configuração do log
type LogOptions struct {
	Level            LogLevel
	Dev              bool
	OutputPaths      []string
	ErrorOutputPaths []string
}

// NewLogger cria um novo logger
func NewLogger(appName string, optionsIn ...LogOptions) (ILogger, error) {
	var (
		cfg         zap.Config
		options     LogOptions
		finalLogger logger
		err         error
	)

	if len(optionsIn) > 0 {
		options = optionsIn[0]
	}

	if options.Dev {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}

	zapLevel, hasZapLevel := convertLevel(options.Level)
	if !hasZapLevel {
		zapLevel = zap.InfoLevel
		options.Level = InfoLevel
	}

	cfg.Encoding = "json"
	cfg.Level = zap.NewAtomicLevelAt(zapLevel)
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.MessageKey = "message"
	cfg.EncoderConfig.LevelKey = "level"
	cfg.EncoderConfig.CallerKey = "caller"
	cfg.EncoderConfig.NameKey = "name"
	cfg.EncoderConfig.TimeKey = "time"
	cfg.EncoderConfig.StacktraceKey = "log_stack_trace"
	cfg.InitialFields = map[string]interface{}{
		"application": appName,
	}

	cfg.OutputPaths = append([]string{"stdout"}, options.OutputPaths...)
	cfg.ErrorOutputPaths = append([]string{"stderr"}, options.ErrorOutputPaths...)

	finalLogger.logger, err = cfg.Build(zap.AddStacktrace(zapLevel))
	if err != nil {
		return nil, err
	}

	finalLogger.level = options.Level
	finalLogger.logger = finalLogger.logger.WithOptions(zap.AddCallerSkip(1))

	return &finalLogger, nil
}

// logger implementa as funções do logger
type logger struct {
	logger *zap.Logger
	level  LogLevel
}

// Level retorna o level do logger
func (l logger) Level() LogLevel {
	return l.level
}

// Info executa um log de level Info
func (l logger) Info(msg string, fields ...LogFields) {
	l.addFields(fields...).logger.Info(msg)
}

// Error executa um log de level Error
func (l logger) Error(erro error, fields ...LogFields) {
	mensagemErro := "Erro detectado"

	// TODO: Adicionar tratamento para error
	fields = append(fields, LogFields{"cause": erro.Error()})

	l.addFields(fields...).logger.Error(mensagemErro)
}

// ErrorMsg executa um log de level Error com mensagem
func (l logger) ErrorMsg(errorMsg string, fields ...LogFields) {
	l.addFields(fields...).logger.Error(errorMsg)
}

// Fatal executa um log de level Fatal
func (l logger) Fatal(msg string, fields ...LogFields) {
	l.addFields(fields...).logger.Fatal(msg)
}

// Panic executa um log de level Panic
func (l logger) Panic(msg string, fields ...LogFields) {
	l.addFields(fields...).logger.Panic(msg)
}

// Debug executa um log de level Debug
func (l logger) Debug(msg string, fields ...LogFields) {
	l.addFields(fields...).logger.Debug(msg)
}

// Warn executa um log de level Warn
func (l logger) Warn(msg string, fields ...LogFields) {
	l.addFields(fields...).logger.Warn(msg)
}

// AddFields adiciona campos ao Logger
func (l logger) AddFields(fields ...LogFields) ILogger {
	return l.addFields(fields...)
}

func (l logger) addFields(fields ...LogFields) logger {
	if len(fields) > 0 {
		for key, value := range fields[0] {
			l.logger = l.logger.With(zap.Any(key, value))
		}
	}

	return l
}

// AddField adiciona um campo ao Logger
func (l logger) AddField(key string, value interface{}) ILogger {
	return l.addFields(LogFields{key: value})
}
