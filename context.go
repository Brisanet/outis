package outis

import (
	"context"
	"errors"
	"time"
)

type period struct {
	startHour, endHour     uint
	startMinute, endMinute uint
}

// Context defines the data structure of the routine context.
type Context struct {
	ID        ID
	RoutineID ID
	Name      string
	Desc      string
	period    period
	Interval  time.Duration
	Path      string
	RunAt     time.Time
	Watcher   Watch

	script            func(*Context) error
	metadata          Metadata
	latency           time.Duration
	notUseLoop        bool
	histogram         []*Histogram
	indicator         []*Indicator
	log               ILogger
	context           context.Context //nolint:containedctx
	contextCancelFunc context.CancelFunc
}

// Context retorna o context.Context.
func (ctx *Context) Context() context.Context {
	return ctx.context
}

// Cancel cancela o contexto.
func (ctx *Context) Cancel() {
	ctx.contextCancelFunc()
}

// Done retorna um canal que espera o canal ser finalizado.
func (ctx *Context) Done() <-chan struct{} {
	return ctx.context.Done()
}

// Err retorna o erro no contexto.
func (ctx *Context) Err() error {
	return ctx.context.Err() //nolint:wrapcheck
}

// Copy cria um cópia do contexto atual.
func (ctx *Context) Copy(baseCtxIn ...context.Context) *Context {
	baseCtx := ctx.context

	if len(baseCtxIn) > 0 {
		baseCtx = baseCtxIn[0]
	}

	childContext, childContextCancelFunc := context.WithCancel(baseCtx)

	return &Context{
		ID:                ctx.ID,
		RoutineID:         ctx.RoutineID,
		Name:              ctx.Name,
		Desc:              ctx.Desc,
		period:            ctx.period,
		Interval:          ctx.Interval,
		Path:              ctx.Path,
		RunAt:             ctx.RunAt,
		Watcher:           ctx.Watcher,
		script:            ctx.script,
		metadata:          ctx.metadata,
		latency:           ctx.latency,
		notUseLoop:        ctx.notUseLoop,
		histogram:         make([]*Histogram, 0),
		indicator:         make([]*Indicator, 0),
		log:               ctx.log,
		context:           childContext,
		contextCancelFunc: childContextCancelFunc,
	}
}

// GetLatency get script execution latency (in seconds).
func (ctx *Context) GetLatency() float64 {
	return ctx.latency.Seconds()
}

// LogInfo executa a função Info do log do contexto.
func (ctx *Context) LogInfo(msg string, fields ...LogFields) {
	ctx.log.Info(msg, fields...)
}

// LogError executa a função Error do log do contexto.
func (ctx *Context) LogError(err error, fields ...LogFields) {
	ctx.log.Error(err, fields...)
}

// LogErrorMsg executa a função ErrorMsg do log do contexto com uma mensagem de erro.
func (ctx *Context) LogErrorMsg(msg string, fields ...LogFields) {
	ctx.log.ErrorMsg(msg, fields...)
}

// LogFatal executa a função Fatal do log do contexto.
func (ctx *Context) LogFatal(msg string, fields ...LogFields) {
	ctx.log.Fatal(msg, fields...)
}

// LogPanic executa a função Panic do log do contexto.
func (ctx *Context) LogPanic(msg string, fields ...LogFields) {
	ctx.log.Panic(msg, fields...)
}

// LogDebug executa a função Debug do log do contexto.
func (ctx *Context) LogDebug(msg string, fields ...LogFields) {
	ctx.log.Debug(msg, fields...)
}

// LogWarn executa a função Warn do log do contexto.
func (ctx *Context) LogWarn(msg string, fields ...LogFields) {
	ctx.log.Warn(msg, fields...)
}

// AddSingleMetadata método adiciona 1 metadata no contexto.
func (ctx *Context) AddSingleMetadata(key string, args interface{}) *Context {
	copyCtx := ctx.Copy()
	copyCtx.metadata.Set(key, args)
	copyCtx.log = copyCtx.log.AddField(key, args)

	return copyCtx
}

// AddMetadata método adiciona metadata no contexto.
func (ctx *Context) AddMetadata(metadata Metadata) *Context {
	copyCtx := ctx.Copy()

	for key, value := range metadata {
		copyCtx.metadata.Set(key, value)
		copyCtx.log = copyCtx.log.AddField(key, value)
	}

	return copyCtx
}

func (ctx *Context) metrics(watch *Watch, now time.Time) {
	watch.outis.Event(ctx, EventMetric{
		ID:         ctx.ID.ToString(),
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
			ID:        ctx.RoutineID.ToString(),
			Name:      ctx.Name,
			Path:      ctx.Path,
			StartedAt: ctx.RunAt,
		},
	})

	ctx.metadata, ctx.indicator, ctx.histogram = Metadata{}, []*Indicator{}, []*Histogram{}
}

func (ctx *Context) sleep(now time.Time) {
	if ctx.mustWait(now.Hour(), ctx.period.startHour, ctx.period.endHour) {
		time.Sleep(time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+int(ctx.period.startHour),
			0, 0, 0, now.Location()).Sub(now))
	}

	if ctx.mustWait(now.Minute(), ctx.period.startMinute, ctx.period.endMinute) {
		time.Sleep(time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+int(ctx.period.startHour),
			now.Minute()+int(ctx.period.startMinute), 0, 0, now.Location()).Sub(now))
	}
}

func (ctx *Context) mustWait(time int, start, end uint) bool {
	if start == 0 && end == 0 {
		return false
	}

	if start <= end {
		return !(time >= int(start) && time <= int(end))
	}

	return !(time >= int(start) || time <= int(end))
}

func (ctx *Context) validate() error {
	if ctx.RoutineID == "" {
		return errors.New("the routine id is required")
	}

	if ctx.Name == "" {
		return errors.New("the routine name is required")
	}

	if ctx.script == nil {
		return errors.New("the routine is required")
	}

	return nil
}
