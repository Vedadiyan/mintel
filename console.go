package mintel

import (
	"fmt"
	"sync"

	"github.com/vedadiyan/mintel/util/template"
)

type (
	ConsoleWriter struct {
		data   map[string]any
		tel    *Console
		binder template.Binder
		mut    sync.RWMutex
	}

	ConsoleLogger struct {
		*ConsoleWriter
	}
	ConsoleTracer struct {
		*ConsoleWriter
	}
	ConsoleMeter struct {
		*ConsoleWriter
	}

	Console struct {
		logger *ConsoleLogger
		tracer *ConsoleTracer
		meter  *ConsoleMeter
		pool   *sync.Pool
	}
)

func NewConsole(l, t, m template.Binder) CreateFunc {
	var pool *sync.Pool
	pool = &sync.Pool{
		New: func() any {
			c := new(Console)
			c.logger = &ConsoleLogger{
				ConsoleWriter: &ConsoleWriter{
					binder: l,
					tel:    c,
					data:   make(map[string]any),
				},
			}
			c.tracer = &ConsoleTracer{
				ConsoleWriter: &ConsoleWriter{
					binder: t,
					tel:    c,
					data:   make(map[string]any),
				},
			}
			c.meter = &ConsoleMeter{
				ConsoleWriter: &ConsoleWriter{
					binder: m,
					tel:    c,
					data:   make(map[string]any),
				},
			}
			c.pool = pool
			return c
		},
	}

	return func(m Metadata) Telemetry {
		tel := pool.Get().(Telemetry)
		return tel
	}
}

func (c *Console) Logger() Writer {
	return c.logger
}

func (c *Console) Tracer() Writer {
	return c.tracer
}

func (c *Console) Meter() Writer {
	return c.meter
}

func (c *Console) Print(v any) {
	fmt.Println(v)
}

func (c *Console) Close() {
	c.Logger().Flush()
	c.Tracer().Flush()
	c.Meter().Flush()
	defer c.pool.Put(c)
	defer c.logger.Clear()
	defer c.tracer.Clear()
	defer c.meter.Clear()
}

func (c *ConsoleWriter) Add(kvs ...*KeyValue) Writer {
	c.mut.Lock()
	defer c.mut.Unlock()
	for _, kv := range kvs {
		c.data[kv.Key] = kv.Value
	}
	return c
}

func (c *ConsoleWriter) Flush() {
	defer c.Clear()
	c.mut.RLock()
	defer c.mut.RUnlock()
	if len(c.data) == 0 {
		return
	}
	c.tel.Print(template.Bind(c.binder, c.data))
}

func (c *ConsoleWriter) Clear() {
	c.mut.Lock()
	defer c.mut.Unlock()
	for key := range c.data {
		delete(c.data, key)
	}
}
