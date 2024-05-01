package outis

import (
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"syscall"
	"time"
)

// Event ...
type Event interface{}

// Routine ...
type Routine func(*Context)

// Option ...
type Option func(*Context)

// ID ...
type ID string

// ToString ...
func (id ID) ToString() string {
	return string(id)
}

// Interface ...
type Interface interface {
	// Lock ...
	Lock(ID) (ID, error)

	// Unlock ...
	Unlock(ID) error

	// Store ...
	Store(*Context) error

	// Load ...
	Load(*Context) error

	// Event ...
	Event(Event)
}

// watch ...
type watch struct {
	name    string
	inter   Interface
	channel chan interface{}
	signal  chan os.Signal
}

// Watcher ...
func Watcher(name string, ioutis Interface) *watch {
	watch := &watch{
		name:    name,
		inter:   ioutis,
		channel: make(chan interface{}),
		signal:  make(chan os.Signal, 1),
	}

	signal.Notify(watch.signal, syscall.SIGINT, syscall.SIGTERM)

	return watch
}

// Wait ...
func (w *watch) Wait() {
	for {
		select {
		case <-w.signal:
			w.inter.Event("closing signal received")
			return
		case err := <-w.channel:
			w.inter.Event(err)
		}
	}
}

// Go ...
func (w *watch) Go(opts ...Option) {
	ctx := &Context{
		channel:      make(chan interface{}),
		interval:     time.Minute,
		metric:       Metric{},
		loadInterval: 0,
	}

	ctx.With(opts...)

	if err := ctx.validate(); err != nil {
		w.inter.Event(err)
		return
	}

	info := runtime.FuncForPC(reflect.ValueOf(ctx.routine).Pointer())
	file, line := info.FileLine(info.Entry())
	ctx.path = fmt.Sprintf("%s:%v", file, line)

	if err := w.inter.Store(ctx); err != nil {
		w.inter.Event(err)
		return
	}

	if ctx.loadInterval != 0 {
		go ctx.reload(w.inter)
	}

	defer func() {
		if r := recover(); r != nil {
			w.inter.Event(fmt.Errorf("PANIC: %v", r))
		}
	}()

	go func() {
		for value := range ctx.channel {
			w.inter.Event(value)
		}
	}()

	ticker := time.NewTicker(ctx.interval)
	for range ticker.C {
		if !ctx.isTime(time.Now().Hour()) {
			continue
		}

		if err := w.process(ctx); err != nil {
			w.inter.Event(err)
			continue
		}

		ticker.Reset(ctx.interval)
	}
}

func (w *watch) process(ctx *Context) error {
	now := time.Now()

	defer func() {
		if r := recover(); r != nil {
			ctx.Error(fmt.Errorf("PANIC: %v", r))
		}
	}()

	id, err := w.inter.Lock(ctx.id)
	if err != nil {
		return err
	}

	w.inter.Event(fmt.Sprintf(`[INITIALIZED] routine '%s' with id '%s'`, ctx.GetName(), id))
	ctx.routine(ctx)

	latency := time.Since(now).Seconds()
	w.inter.Event(Metrics{
		ID:          id.ToString(),
		Initialized: now,
		Terminated:  time.Now(),
		Latency:     latency,
		Metadata:    ctx.metric,
		Routine: RoutineMetric{
			ID:   ctx.GetID().ToString(),
			Name: ctx.GetName(),
			Path: ctx.path,
		},
	})

	ctx.metric = Metric{}
	if err = w.inter.Unlock(id); err != nil {
		return err
	}

	w.inter.Event(fmt.Sprintf(`[TERMINATED] routine '%s' with id '%s' in %v seconds`, ctx.GetName(), id, latency))
	return nil
}

func (ctx *Context) reload(ioutis Interface) {
	ticker := time.NewTicker(ctx.loadInterval)
	for range ticker.C {
		if err := ioutis.Load(ctx); err != nil {
			ioutis.Event(err)
		}
		ticker.Reset(ctx.loadInterval)
		ioutis.Event(fmt.Sprintf(`[UPDATED] routine '%s' with id '%s' has been updated`, ctx.GetName(), ctx.GetID()))
	}
}
