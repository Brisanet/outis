package outis

// Option defines the option type of a routine
type Option func(*ContextImpl)

// WithName defines the name of a routine
func WithName(name string) Option {
	return func(ctx *ContextImpl) { ctx.name = name }
}

// WithDesc defines the description of a routine
func WithDesc(desc string) Option {
	return func(ctx *ContextImpl) { ctx.Desc = desc }
}

// WithID defines a routine's identifier
func WithID(id ID) Option {
	return func(ctx *ContextImpl) { ctx.routineID = id }
}

// WithScript defines the script function that will be executed
func WithScript(fn func(Context) error) Option {
	return func(ctx *ContextImpl) { ctx.script = fn }
}

// WithInterval defines the interval at which the script will be executed
func WithInterval(interval Interval) Option {
	return func(ctx *ContextImpl) { ctx.Interval = &interval }
}

// WithNotUseLoop define that the routine will not enter a loop
func WithNotUseLoop() Option {
	return func(ctx *ContextImpl) { ctx.notUseLoop = true }
}

// WithExecuteFirstTimeNow define that the routine will execute first time when Watcher.Go is called
func WithExecuteFirstTimeBeforeInterval() Option {
	return func(ctx *ContextImpl) { ctx.executeFirstTimeBeforeInterval = true }
}

// WatcherOption defines the option type of a watcher
type WatcherOption func(*Watch)

// Logger defines the implementation of the log interface
func Logger(logger ILogger) WatcherOption {
	return func(watch *Watch) { watch.log = logger }
}

// Impl defines the implementation of the main interface
func Impl(outis IOutis) WatcherOption {
	return func(watch *Watch) { watch.outis = outis }
}
