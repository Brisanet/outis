package outis

// Event sets the type for the event structure
type Event interface{}

// IOutis is the main interface for implementing the outis lib.
type IOutis interface {
	Go(fn func() error)
	Wait() error

	Init(ctx *Context) error
	Before(ctx *Context) error
	After(ctx *Context) error
	Event(ctx *Context, event Event)
}

// ILogger methods for logging messages.
type ILogger interface {
	AddFields(fields ...Metadata) ILogger
	AddField(key string, value interface{}) ILogger
	Info(msg string, fields ...Metadata)
	Error(erro error, fields ...Metadata)
	ErrorMsg(errorMsg string, fields ...Metadata)
	Fatal(msg string, fields ...Metadata)
	Panic(msg string, fields ...Metadata)
	Debug(msg string, fields ...Metadata)
	Warn(msg string, fields ...Metadata)
}
