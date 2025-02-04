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

// Context defines the data structure of the routine context
type Context struct {
	Id        ID
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
	histrogram        []*histogram
	indicator         []*indicator
	log               ILogger
	context           context.Context
	contextCancelFunc context.CancelFunc
}

func (ctx *Context) SetContext(newCtx context.Context) *Context {
	copy := ctx.copy()
	copy.context = newCtx

	return copy
}

// Context returns the application context
func (ctx *Context) Context() context.Context {
	return ctx.context
}

func (ctx *Context) Cancel() {
	ctx.contextCancelFunc()
}

func (ctx *Context) Done() <-chan struct{} {
	return ctx.context.Done()
}

func (ctx *Context) Err() error {
	return ctx.context.Err()
}

func (ctx *Context) copy() *Context {
	childContext, childContextCancelFunc := context.WithCancel(ctx.context)

	return &Context{
		indicator:         ctx.indicator,
		metadata:          ctx.metadata,
		log:               ctx.log,
		Interval:          ctx.Interval,
		RunAt:             ctx.RunAt,
		Watcher:           ctx.Watcher,
		context:           childContext,
		contextCancelFunc: childContextCancelFunc,
	}
}

// GetLatency get script execution latency (in seconds)
func (ctx *Context) GetLatency() float64 {
	return ctx.latency.Seconds()
}

// Info executa a função Info do log do contexto
func (ctx *Context) LogInfo(msg string, fields ...LogFields) {
	ctx.log.Info(msg, fields...)
}

// Error executa a função Erro do log do contexto
func (ctx *Context) LogError(err error, fields ...LogFields) {
	ctx.log.Error(err, fields...)
}

// ErrorMsg executa a função Erro do log do contexto com uma mensagem de erro
func (ctx *Context) LogErrorMsg(msg string, fields ...LogFields) {
	ctx.log.ErrorMsg(msg, fields...)
}

func (ctx *Context) LogFatal(msg string, fields ...LogFields) {
	ctx.log.Fatal(msg, fields...)
}

func (ctx *Context) LogPanic(msg string, fields ...LogFields) {
	ctx.log.Panic(msg, fields...)
}

func (ctx *Context) LogDebug(msg string, fields ...LogFields) {
	ctx.log.Debug(msg, fields...)
}

func (ctx *Context) LogWarn(msg string, fields ...LogFields) {
	ctx.log.Warn(msg, fields...)
}

// Metadata method for adding data to routine metadata
func (ctx *Context) AddSingleMetadata(key string, args interface{}) *Context {
	copy := ctx.copy()
	copy.metadata.Set(key, args)
	copy.log = copy.log.AddField(key, args)

	return copy
}

// Metadata method for adding data to routine metadata
func (ctx *Context) AddMetadata(metadata Metadata) *Context {
	copy := ctx.copy()

	for key, value := range metadata {
		copy.metadata.Set(key, value)
		copy.log = copy.log.AddField(key, value)
	}

	return copy
}

func (ctx *Context) metrics(w *Watch, now time.Time) {
	w.outis.Event(ctx, EventMetric{
		ID:         ctx.Id.ToString(),
		StartedAt:  now,
		FinishedAt: time.Now(),
		Latency:    time.Since(now),
		Metadata:   ctx.metadata,
		Indicators: ctx.indicator,
		Histograms: ctx.histrogram,
		Watcher: WatcherMetric{
			ID:    w.Id.ToString(),
			Name:  w.Name,
			RunAt: w.RunAt,
		},
		Routine: RoutineMetric{
			ID:        ctx.RoutineID.ToString(),
			Name:      ctx.Name,
			Path:      ctx.Path,
			StartedAt: ctx.RunAt,
		},
	})

	ctx.metadata, ctx.indicator, ctx.histrogram = Metadata{}, []*indicator{}, []*histogram{}
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
