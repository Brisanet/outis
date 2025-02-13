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
	Level() LogLevel
	Info(msg string, fields ...LogFields)
	Error(erro error, fields ...LogFields)
	ErrorMsg(errorMsg string, fields ...LogFields)
	Fatal(msg string, fields ...LogFields)
	Panic(msg string, fields ...LogFields)
	Debug(msg string, fields ...LogFields)
	Warn(msg string, fields ...LogFields)
	AddFields(fields ...LogFields) ILogger
	AddField(key string, value interface{}) ILogger
}
