package outis

import (
	"fmt"
	"math/rand"
	"strconv"

	"golang.org/x/sync/errgroup"
)

type server struct {
	errGroup errgroup.Group
}

func newOutis() IOutis {
	return &server{
		errGroup: errgroup.Group{},
	}
}

// Go executa a função passada
func (s *server) Go(fn func() error) {
	s.errGroup.Go(fn)
}

// Wait espera a execução da função Go
func (s *server) Wait() error {
	return s.errGroup.Wait()
}

// Init implements a business rule when initializing a routine
func (s *server) Init(ctx *Context) error {
	ctx.LogInfo(fmt.Sprintf("script '%s' (rid: %s) initialized", ctx.Name, ctx.RoutineID))
	return nil
}

// Before implements a business rule before initializing script execution
func (s *server) Before(ctx *Context) error {
	ctx.ID = ID(strconv.FormatInt(rand.Int63(), 10))
	ctx.LogInfo(fmt.Sprintf("script '%s' (rid: %s, id: %s) initialized", ctx.Name, ctx.RoutineID, ctx.ID))
	return nil
}

// After implements a business rule after initializing script execution
func (s *server) After(ctx *Context) error {
	ctx.LogInfo(fmt.Sprintf("script '%s' (rid: %s, id: %s) finished", ctx.Name, ctx.RoutineID, ctx.ID))
	return nil
}

// Event implements a business rule for event handling
func (s *server) Event(ctx *Context, event Event) {
	if metric, ok := event.(EventMetric); ok {
		ctx.LogDebug("Metrics", LogFields{"metrics": metric})
	}
}
