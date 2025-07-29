package outis

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

// ID defines the type of identifier
type ID string

// ToString return the identifier as a string
func (id ID) ToString() string {
	return string(id)
}

// Watch defines the type of the watcher structure
type Watch struct {
	Id    ID        `json:"id"`
	Name  string    `json:"name"`
	RunAt time.Time `json:"run_at"`

	outis IOutis
	log   ILogger
}

// Watcher initializes a new watcher
func Watcher(id, name string, opts ...WatcherOption) *Watch {
	watch := &Watch{
		Id:    ID(id),
		Name:  name,
		outis: newOutis(),
		RunAt: time.Now(),
	}

	for _, opt := range opts {
		opt(watch)
	}
	if watch.log == nil {
		logger, err := NewLogger(name)
		if err != nil {
			log.Fatal(err)
		}
		watch.log = logger
	}

	return watch
}

// Wait method responsible for keeping routines running
func (watch *Watch) Wait() {
	if err := watch.outis.Wait(); err != nil {
		watch.log.Error(err)
		return
	}
}

// Wait responsible for keeping routines running
func Wait() {
	wait := make(chan os.Signal, 1)
	signal.Notify(wait, syscall.SIGINT, syscall.SIGTERM)

	for range wait {
		return
	}
}

// Go create a new routine in the watcher
func (watch *Watch) Go(opts ...Option) {
	watch.outis.Go(func() error {
		var (
			childContext, childContextCancelFunc = context.WithCancel(context.Background())
			ctx                                  = &ContextImpl{
				id:                ID(strconv.FormatInt(rand.Int63(), 10)),
				indicator:         make([]*Indicator, 0),
				metadata:          make(Metadata),
				log:               watch.log,
				Interval:          time.Minute,
				RunAt:             time.Now(),
				Watcher:           *watch,
				context:           childContext,
				contextCancelFunc: childContextCancelFunc,
			}
			err error
		)

		for _, opt := range opts {
			opt(ctx)
		}

		if err = ctx.validate(); err != nil {
			return err
		}

		info := runtime.FuncForPC(reflect.ValueOf(ctx.script).Pointer())
		file, line := info.FileLine(info.Entry())
		ctx.Path = fmt.Sprintf("%s:%v", file, line)

		if err := watch.outis.Init(ctx); err != nil {
			return err
		}

		defer func() {
			if r := recover(); r != nil {
				ctx.log.Panic(fmt.Sprintf("%v", r))
			}
		}()

		// TODO: refactor the execution logic below when add test

		if ctx.notUseLoop {
			ctx.sleep(time.Now())
			return ctx.execute()
		}

		if ctx.executeFirstTimeBeforeInterval {
			ctx.sleep(time.Now())
			if err = ctx.execute(); err != nil {
				ctx.log.Error(err)
			}
		}

		ticker := time.NewTicker(ctx.Interval)
		for {
			ctx.sleep(time.Now())

			ctx.log.Info("Starting ticker " + ctx.Interval.String())
			ticker.Reset(ctx.Interval)
			select {
			// Espera o contexto ser finalizado
			case <-ctx.context.Done():
				return ctx.context.Err()
			// Espera a próxima execução com base no ticker
			case _, isOpen := <-ticker.C:
				ticker.Stop()
				if !isOpen {
					return nil
				}

				var err error

				// Caso o contexto esteja com erro o script é finalizado
				if err = ctx.context.Err(); err != nil {
					return err
				}

				if err = ctx.execute(); err != nil {
					ctx.log.Error(err)
					continue
				}
			}
		}
	})
}

func (ctx *ContextImpl) execute() error {
	initialTime := time.Now()
	defer func() {
		if r := recover(); r != nil {
			ctx.log.Panic(fmt.Sprintf("%v", r))
		}
	}()

	if err := ctx.Watcher.outis.Before(ctx); err != nil {
		return err
	}

	if err := ctx.script(ctx.Copy()); err != nil {
		return err
	}

	ctx.latency = time.Since(initialTime)
	if err := ctx.Watcher.outis.After(ctx); err != nil {
		return err
	}

	ctx.metrics(&ctx.Watcher, initialTime)

	return nil
}
