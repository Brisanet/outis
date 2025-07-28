package outis

import (
	"context"
	"errors"
	"time"
)

type period struct {
	hourSet                bool
	startHour, endHour     uint
	minuteSet              bool
	startMinute, endMinute uint
}

// Context defines the data structure of the routine context.
type Context interface {
	Context() context.Context
	Cancel()
	Done() <-chan struct{}
	Err() error
	Copy(baseCtxIn ...context.Context) Context
	GetLatency() float64
	LogInfo(msg string, fields ...LogFields)
	LogError(err error, fields ...LogFields)
	LogErrorMsg(msg string, fields ...LogFields)
	LogFatal(msg string, fields ...LogFields)
	LogPanic(msg string, fields ...LogFields)
	LogDebug(msg string, fields ...LogFields)
	LogWarn(msg string, fields ...LogFields)
	AddSingleMetadata(key string, args interface{}) Context
	AddMetadata(metadata Metadata) Context

	Name() string
	RoutineID() ID
	ID() ID
}

// ContextImpl implements context interface
type ContextImpl struct {
	id        ID
	routineID ID
	name      string
	Desc      string
	period    period
	Interval  *Interval
	Path      string
	RunAt     time.Time
	Watcher   Watch

	// Reminder: If new fields are added, change context.Copy function accordingly.
	script                         func(Context) error
	metadata                       Metadata
	latency                        time.Duration
	notUseLoop                     bool
	executeFirstTimeBeforeInterval bool
	histogram                      []*Histogram
	indicator                      []*Indicator
	log                            ILogger
	context                        context.Context //nolint:containedctx
	contextCancelFunc              context.CancelFunc
}

// Context retorna o context.Context.
func (ctx *ContextImpl) Context() context.Context {
	return ctx.context
}

// Cancel cancela o contexto.
func (ctx *ContextImpl) Cancel() {
	ctx.contextCancelFunc()
}

// Done retorna um canal que espera o canal ser finalizado.
func (ctx *ContextImpl) Done() <-chan struct{} {
	return ctx.context.Done()
}

// Err retorna o erro no contexto.
func (ctx *ContextImpl) Err() error {
	return ctx.context.Err() //nolint:wrapcheck
}

// Copy cria um cópia do contexto atual.
func (ctx *ContextImpl) Copy(baseCtxIn ...context.Context) Context {
	return ctx.copy(baseCtxIn...)
}

// Copy cria um cópia do contexto atual.
func (ctx *ContextImpl) copy(baseCtxIn ...context.Context) *ContextImpl {
	baseCtx := ctx.context

	if len(baseCtxIn) > 0 {
		baseCtx = baseCtxIn[0]
	}

	childContext, childContextCancelFunc := context.WithCancel(baseCtx)

	return &ContextImpl{
		id:                             ctx.id,
		routineID:                      ctx.routineID,
		name:                           ctx.name,
		Desc:                           ctx.Desc,
		period:                         ctx.period,
		Interval:                       ctx.Interval,
		Path:                           ctx.Path,
		RunAt:                          ctx.RunAt,
		Watcher:                        ctx.Watcher,
		script:                         ctx.script,
		metadata:                       ctx.metadata,
		latency:                        ctx.latency,
		notUseLoop:                     ctx.notUseLoop,
		executeFirstTimeBeforeInterval: ctx.executeFirstTimeBeforeInterval,
		histogram:                      make([]*Histogram, 0),
		indicator:                      make([]*Indicator, 0),
		log:                            ctx.log,
		context:                        childContext,
		contextCancelFunc:              childContextCancelFunc,
	}
}

// GetLatency get script execution latency (in seconds).
func (ctx *ContextImpl) GetLatency() float64 {
	return ctx.latency.Seconds()
}

// LogInfo executa a função Info do log do contexto.
func (ctx *ContextImpl) LogInfo(msg string, fields ...LogFields) {
	ctx.log.Info(msg, fields...)
}

// LogError executa a função Error do log do contexto.
func (ctx *ContextImpl) LogError(err error, fields ...LogFields) {
	ctx.log.Error(err, fields...)
}

// LogErrorMsg executa a função ErrorMsg do log do contexto com uma mensagem de erro.
func (ctx *ContextImpl) LogErrorMsg(msg string, fields ...LogFields) {
	ctx.log.ErrorMsg(msg, fields...)
}

// LogFatal executa a função Fatal do log do contexto.
func (ctx *ContextImpl) LogFatal(msg string, fields ...LogFields) {
	ctx.log.Fatal(msg, fields...)
}

// LogPanic executa a função Panic do log do contexto.
func (ctx *ContextImpl) LogPanic(msg string, fields ...LogFields) {
	ctx.log.Panic(msg, fields...)
}

// LogDebug executa a função Debug do log do contexto.
func (ctx *ContextImpl) LogDebug(msg string, fields ...LogFields) {
	ctx.log.Debug(msg, fields...)
}

// LogWarn executa a função Warn do log do contexto.
func (ctx *ContextImpl) LogWarn(msg string, fields ...LogFields) {
	ctx.log.Warn(msg, fields...)
}

// AddSingleMetadata método adiciona 1 metadata no contexto.
func (ctx *ContextImpl) AddSingleMetadata(key string, args interface{}) Context {
	copyCtx := ctx.copy()
	copyCtx.metadata.Set(key, args)
	copyCtx.log = copyCtx.log.AddField(key, args)

	return copyCtx
}

// AddMetadata método adiciona metadata no contexto.
func (ctx *ContextImpl) AddMetadata(metadata Metadata) Context {
	copyCtx := ctx.copy()

	for key, value := range metadata {
		copyCtx.metadata.Set(key, value)
		copyCtx.log = copyCtx.log.AddField(key, value)
	}

	return copyCtx
}

func (ctx *ContextImpl) metrics(watch *Watch, now time.Time) {
	watch.outis.Event(ctx, EventMetric{
		ID:         ctx.id.ToString(),
		StartedAt:  now,
		FinishedAt: time.Now(),
		Latency:    time.Since(now),
		Metadata:   ctx.metadata,
		Indicators: ctx.indicator,
		Histograms: ctx.histogram,
		Watcher: WatcherMetric{
			ID:    watch.Id.ToString(),
			Name:  watch.Name,
			RunAt: watch.RunAt,
		},
		Routine: RoutineMetric{
			ID:        ctx.routineID.ToString(),
			Name:      ctx.name,
			Path:      ctx.Path,
			StartedAt: ctx.RunAt,
		},
	})

	ctx.metadata, ctx.indicator, ctx.histogram = Metadata{}, []*Indicator{}, []*Histogram{}
}

func (ctx *ContextImpl) validate() error {
	if ctx.RoutineID() == "" {
		return errors.New("the routine id is required")
	}

	if ctx.name == "" {
		return errors.New("the routine name is required")
	}

	if ctx.script == nil {
		return errors.New("the routine is required")
	}

	if err := ctx.Interval.validate(); err != nil {
		return err
	}

	return nil
}

// Name returns the name of the routine
func (ctx *ContextImpl) Name() string {
	return ctx.name
}

// RoutineID returns the ID of the routine
func (ctx *ContextImpl) RoutineID() ID {
	return ctx.routineID
}

// ID returns the execution ID of the routine
func (ctx *ContextImpl) ID() ID {
	return ctx.id
}

func (ctx *ContextImpl) WaitNextExecution(isFirstScriptExecution bool, executeChan chan struct{}) {
	if ctx.Interval != nil {
		ctx.Interval.Wait(context.Background(), time.Now(), isFirstScriptExecution)
	}
	executeChan <- struct{}{}
}
