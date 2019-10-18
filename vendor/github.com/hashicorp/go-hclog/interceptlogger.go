package hclog

import (
	"sync"
	"sync/atomic"
)

var _ Logger = &interceptLogger{}

type interceptLogger struct {
	Logger

	sync.Mutex
	sinkCount *int32
	Sinks     map[SinkAdapter]struct{}
}

func NewInterceptLogger(opts *LoggerOptions) InterceptLogger {
	intercept := &interceptLogger{
		Logger:    New(opts),
		sinkCount: new(int32),
		Sinks:     make(map[SinkAdapter]struct{}),
	}

	atomic.StoreInt32(intercept.sinkCount, 0)

	return intercept
}

// Emit the message and args at TRACE level to log and sinks
func (i *interceptLogger) Trace(msg string, args ...interface{}) {
	i.Logger.Trace(msg, args...)
	if atomic.LoadInt32(i.sinkCount) == 0 {
		return
	}

	i.Lock()
	defer i.Unlock()
	for s := range i.Sinks {
		s.Accept(i.Name(), Trace, msg, i.retrieveImplied(args...)...)
	}
}

// Emit the message and args at DEBUG level to log and sinks
func (i *interceptLogger) Debug(msg string, args ...interface{}) {
	i.Logger.Debug(msg, args...)
	if atomic.LoadInt32(i.sinkCount) == 0 {
		return
	}

	i.Lock()
	defer i.Unlock()
	for s := range i.Sinks {
		s.Accept(i.Name(), Debug, msg, i.retrieveImplied(args...)...)
	}
}

// Emit the message and args at INFO level to log and sinks
func (i *interceptLogger) Info(msg string, args ...interface{}) {
	i.Logger.Info(msg, args...)
	if atomic.LoadInt32(i.sinkCount) == 0 {
		return
	}

	i.Lock()
	defer i.Unlock()
	for s := range i.Sinks {
		s.Accept(i.Name(), Info, msg, i.retrieveImplied(args...)...)
	}
}

// Emit the message and args at WARN level to log and sinks
func (i *interceptLogger) Warn(msg string, args ...interface{}) {
	i.Logger.Warn(msg, args...)
	if atomic.LoadInt32(i.sinkCount) == 0 {
		return
	}

	i.Lock()
	defer i.Unlock()
	for s := range i.Sinks {
		s.Accept(i.Name(), Warn, msg, i.retrieveImplied(args...)...)
	}
}

// Emit the message and args at ERROR level to log and sinks
func (i *interceptLogger) Error(msg string, args ...interface{}) {
	i.Logger.Error(msg, args...)
	if atomic.LoadInt32(i.sinkCount) == 0 {
		return
	}

	i.Lock()
	defer i.Unlock()
	for s := range i.Sinks {
		s.Accept(i.Name(), Error, msg, i.retrieveImplied(args...)...)
	}
}

func (i *interceptLogger) retrieveImplied(args ...interface{}) []interface{} {
	top := i.Logger.ImpliedArgs()

	cp := make([]interface{}, len(top)+len(args))
	copy(cp, top)
	copy(cp[len(top):], args)

	return cp
}

// Create a new sub-Logger that a name decending from the current name.
// This is used to create a subsystem specific Logger.
// Registered sinks will subscribe to these messages as well.
func (i *interceptLogger) Named(name string) Logger {
	var sub interceptLogger

	sub = *i

	sub.Logger = i.Logger.Named(name)

	return &sub
}

// Create a new sub-Logger with an explicit name. This ignores the current
// name. This is used to create a standalone logger that doesn't fall
// within the normal hierarchy. Registered sinks will subscribe
// to these messages as well.
func (i *interceptLogger) ResetNamed(name string) Logger {
	var sub interceptLogger

	sub = *i

	sub.Logger = i.Logger.ResetNamed(name)

	return &sub
}

// Create a new sub-Logger that a name decending from the current name.
// This is used to create a subsystem specific Logger.
// Registered sinks will subscribe to these messages as well.
func (i *interceptLogger) NamedIntercept(name string) InterceptLogger {
	var sub interceptLogger

	sub = *i

	sub.Logger = i.Logger.Named(name)

	return &sub
}

// Create a new sub-Logger with an explicit name. This ignores the current
// name. This is used to create a standalone logger that doesn't fall
// within the normal hierarchy. Registered sinks will subscribe
// to these messages as well.
func (i *interceptLogger) ResetNamedIntercept(name string) InterceptLogger {
	var sub interceptLogger

	sub = *i

	sub.Logger = i.Logger.ResetNamed(name)

	return &sub
}

// Return a sub-Logger for which every emitted log message will contain
// the given key/value pairs. This is used to create a context specific
// Logger.
func (i *interceptLogger) With(args ...interface{}) Logger {
	var sub interceptLogger

	sub = *i

	sub.Logger = i.Logger.With(args...)

	return &sub
}

// RegisterSink attaches a SinkAdapter to interceptLoggers sinks.
func (i *interceptLogger) RegisterSink(sink SinkAdapter) {
	i.Lock()
	defer i.Unlock()

	i.Sinks[sink] = struct{}{}

	atomic.AddInt32(i.sinkCount, 1)
}

// DeregisterSink removes a SinkAdapter from interceptLoggers sinks.
func (i *interceptLogger) DeregisterSink(sink SinkAdapter) {
	i.Lock()
	defer i.Unlock()

	delete(i.Sinks, sink)

	atomic.AddInt32(i.sinkCount, -1)
}
