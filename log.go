package outis

import (
	"runtime"
	"strconv"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogFields map[string]interface{}

type LogLevel string

const (
	// DebugLevel representa o nível de debug do Logger
	DebugLevel LogLevel = "DebugLevel"
	// InfoLevel representa o nível de info do Logger
	InfoLevel LogLevel = "InfoLevel"
	// WarnLevel representa o nível de warn do Logger
	WarnLevel LogLevel = "WarnLevel"
	// ErrorLevel representa o nível de error do Logger
	ErrorLevel LogLevel = "ErrorLevel"
	// DPanicLevel representa o nível de dpanic do Logger
	DPanicLevel LogLevel = "DPanicLevel"
	// PanicLevel representa o nível de panic do Logger
	PanicLevel LogLevel = "PanicLevel"
	// FatalLevel representa o nível de fatal do Logger
	FatalLevel LogLevel = "FatalLevel"
)

type Options struct {
	Level LogLevel
	Dev   bool
}

func NewLogger(appName string, optionsIn ...Options) (ILogger, error) {
	var (
		cfg         zap.Config
		logLevelZap zap.AtomicLevel
		logLevel    LogLevel
		options     Options
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

	switch options.Level {
	case DebugLevel:
		logLevelZap = zap.NewAtomicLevelAt(zap.DebugLevel)
		logLevel = DebugLevel
	default:
		logLevelZap = zap.NewAtomicLevelAt(zap.InfoLevel)
		logLevel = InfoLevel
	}

	cfg.Encoding = "json"
	cfg.Level = logLevelZap
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

	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}

	finalLogger.logger, err = cfg.Build()
	if err != nil {
		return nil, err
	}

	finalLogger.level = logLevel
	finalLogger.logger = finalLogger.logger.WithOptions(zap.AddCallerSkip(1))

	return &finalLogger, nil
}

// logger implementa as funções do logger
type logger struct {
	logger *zap.Logger
	level  LogLevel
}

func (l logger) Level() LogLevel {
	return l.level
}

func (l logger) Info(msg string, fields ...LogFields) {
	l.addStackTrace().addFields(fields...).logger.Info(msg)
}

func (l logger) Error(erro error, fields ...LogFields) {
	var (
		mensagemErro = "Erro detectado"
	)

	// TODO: Adicionar tratamento para error
	fields = append(fields, LogFields{"cause": erro.Error()})

	l.addStackTrace().addFields(fields...).logger.Error(mensagemErro)
}

func (l logger) ErrorMsg(errorMsg string, fields ...LogFields) {
	l.addStackTrace().addFields(fields...).logger.Error(errorMsg)
}

func (l logger) Fatal(msg string, fields ...LogFields) {
	l.addStackTrace().addFields(fields...).logger.Fatal(msg)
}

func (l logger) Panic(msg string, fields ...LogFields) {
	l.addStackTrace().addFields(fields...).logger.Panic(msg)
}

func (l logger) Debug(msg string, fields ...LogFields) {
	l.addStackTrace().addFields(fields...).logger.Debug(msg)
}

func (l logger) Warn(msg string, fields ...LogFields) {
	l.addStackTrace().addFields(fields...).logger.Warn(msg)
}

// AddFields adiciona campos ao Logger
func (l logger) AddFields(fields ...LogFields) ILogger {
	return l.addFields(fields...)
}

// AddFields adiciona campos ao Logger
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

func (l logger) addStackTrace() logger {
	return l.addFields(LogFields{"stack_trace": getStackTrace()})
}

func getStackTrace() (stackTrace []string) {
	const depth = 32
	var pcs [depth]uintptr
	runtime.Callers(3, pcs[:])

	for ii := range pcs {
		pc := pcs[ii] - 1
		functionFound := runtime.FuncForPC(pc)
		if functionFound == nil {
			continue
		}
		file, line := functionFound.FileLine(pc)
		stackTrace = append(stackTrace, file+":"+strconv.Itoa(line))
	}

	return stackTrace
}
