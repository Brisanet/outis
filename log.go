package outis

import (
	"runtime"
	"strconv"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level string

const (
	// DebugLevel representa o nível de debug do Logger
	DebugLevel Level = "DebugLevel"
	// InfoLevel representa o nível de info do Logger
	InfoLevel Level = "InfoLevel"
	// WarnLevel representa o nível de warn do Logger
	WarnLevel Level = "WarnLevel"
	// ErrorLevel representa o nível de error do Logger
	ErrorLevel Level = "ErrorLevel"
	// DPanicLevel representa o nível de dpanic do Logger
	DPanicLevel Level = "DPanicLevel"
	// PanicLevel representa o nível de panic do Logger
	PanicLevel Level = "PanicLevel"
	// FatalLevel representa o nível de fatal do Logger
	FatalLevel Level = "FatalLevel"
)

type Options struct {
	Level Level
	Prod  bool
}

func NewLogger(appName string, optionsIn ...Options) (ILogger, error) {
	var (
		cfg         zap.Config
		logLevelZap zap.AtomicLevel
		options     Options
		finalLogger logger
		err         error
	)

	if len(optionsIn) > 0 {
		options = optionsIn[0]
	}

	if options.Prod {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
	}

	switch options.Level {
	case DebugLevel:
		logLevelZap = zap.NewAtomicLevelAt(zap.DebugLevel)
	default:
		logLevelZap = zap.NewAtomicLevelAt(zap.InfoLevel)
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

	finalLogger.logger = finalLogger.logger.WithOptions(zap.AddCallerSkip(1))

	return &finalLogger, nil
}

// logger implementa as funções do logger
type logger struct {
	logger *zap.Logger
}

func (l logger) Info(msg string, fields ...Metadata) {
	l.addStackTrace().addFields(fields...).logger.Info(msg)
}

func (l logger) Error(erro error, fields ...Metadata) {
	var (
		mensagemErro = "Erro detectado"
	)

	// TODO: Adicionar tratamento para error
	l = l.addStackTrace()
	fields = append(fields, Metadata{"erro": erro.Error()})

	l.addFields(fields...).logger.Error(mensagemErro)
}

func (l logger) ErrorMsg(errorMsg string, fields ...Metadata) {
	l.addStackTrace().addFields(fields...).logger.Error(errorMsg)
}

func (l logger) Fatal(msg string, fields ...Metadata) {
	l.addStackTrace().addFields(fields...).logger.Fatal(msg)
}

func (l logger) Panic(msg string, fields ...Metadata) {
	l.addStackTrace().addFields(fields...).logger.Panic(msg)
}

func (l logger) Debug(msg string, fields ...Metadata) {
	l.addStackTrace().addFields(fields...).logger.Debug(msg)
}

func (l logger) Warn(msg string, fields ...Metadata) {
	l.addStackTrace().addFields(fields...).logger.Warn(msg)
}

// AddFields adiciona campos ao Logger
func (l logger) AddFields(fields ...Metadata) ILogger {
	return l.addFields(fields...)
}

// AddFields adiciona campos ao Logger
func (l logger) addFields(fields ...Metadata) logger {
	if len(fields) > 0 {
		for key, value := range fields[0] {
			l.logger = l.logger.With(zap.Any(key, value))
		}
	}

	return l
}

// AddField adiciona um campo ao Logger
func (l logger) AddField(key string, value interface{}) ILogger {
	return l.addField(key, value)
}

// AddField adiciona um campo ao Logger
func (l logger) addField(key string, value interface{}) logger {
	return l.addFields(Metadata{key: value})
}

func (l logger) addStackTrace() logger {
	return l.addField("stack_trace", getStackTrace())
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
