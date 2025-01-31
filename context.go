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

	script     func(*Context) error
	metadata   Metadata
	latency    time.Duration
	notUseLoop bool
	histrogram []*histogram
	indicator  []*indicator
	log        ILogger
	context    context.Context
	// TODO: cancelCtx
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

// TODO: context.Context -> outis.Context

func (ctx *Context) copy() *Context {
	parentContext, _ := context.WithCancel(ctx.context) //nolint

	return &Context{
		indicator: ctx.indicator,
		metadata:  ctx.metadata,
		log:       ctx.log,
		Interval:  ctx.Interval,
		RunAt:     ctx.RunAt,
		Watcher:   ctx.Watcher,
		context:   parentContext,
	}
}

// GetLatency get script execution latency (in seconds)
func (ctx *Context) GetLatency() float64 {
	return ctx.latency.Seconds()
}

// Info executa a função Info do log do contexto
func (ctx *Context) Info(msg string, fields ...Metadata) {
	ctx.log.Info(msg, fields...)
}

func (ctx *Context) Debug(msg string, fields ...Metadata) {
	ctx.log.Debug(msg, fields...)
}

// Error executa a função Erro do log do contexto
func (ctx *Context) Error(err error, fields ...Metadata) {
	ctx.log.Error(err, fields...)
}

// ErrorMsg executa a função Erro do log do contexto com uma mensagem de erro
func (ctx *Context) ErrorMsg(msg string, fields ...Metadata) {
	ctx.log.ErrorMsg(msg, fields...)
}

func (ctx *Context) Warn(msg string, fields ...Metadata) {
	ctx.log.Warn(msg, fields...)
}

func (ctx *Context) Panic(msg string, fields ...Metadata) {
	ctx.log.Panic(msg, fields...)
}

func (ctx *Context) Fatal(msg string, fields ...Metadata) {
	ctx.log.Fatal(msg, fields...)
}

// Metadata method for adding data to routine metadata
func (ctx *Context) AddSingleMetadata(key string, args interface{}) *Context {
	copy := ctx.copy()
	copy.metadata.Set(key, args)
	copy.log.AddField(key, args)

	return copy
}

// Metadata method for adding data to routine metadata
func (ctx *Context) AddMetadata(metadata Metadata) *Context {
	copy := ctx.copy()

	for key, value := range metadata {
		copy.metadata.Set(key, value)
		copy.log.AddField(key, value)
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
