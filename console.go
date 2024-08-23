package mintel

import (
	"log"

	"github.com/vedadiyan/mintel/util/template"
)

type (
	ConsoleLogger struct {
		binder template.Binder
	}
	ConsoleTracer struct {
		data   map[string]any
		binder template.Binder
	}
	ConsoleMeter struct {
		data   map[string]any
		binder template.Binder
	}
	ConsoleClient struct {
		logger Logger
		tracer Tracer
		meter  Meter
	}
)

func NewConsole(logTmpl, traceTmpl, metricsTmpl string) (*ConsoleClient, error) {
	l, err := template.Parse(logTmpl)
	if err != nil {
		return nil, err
	}
	t, err := template.Parse(traceTmpl)
	if err != nil {
		return nil, err
	}
	m, err := template.Parse(metricsTmpl)
	if err != nil {
		return nil, err
	}
	c := new(ConsoleClient)
	c.logger = &ConsoleLogger{
		binder: l,
	}
	c.tracer = &ConsoleTracer{
		binder: t,
		data:   make(map[string]any),
	}
	c.meter = &ConsoleMeter{
		binder: m,
		data:   make(map[string]any),
	}
	return c, nil
}

func (c *ConsoleClient) Logger() Logger {
	return c.logger
}

func (c *ConsoleClient) Tracer() Tracer {
	return c.tracer
}

func (c *ConsoleClient) Meter() Meter {
	return c.meter
}

func (c *ConsoleClient) Close(v ...any) {
	c.Tracer().Notify()
	c.Meter().Notify()
	log.Println("DONE", v)
}

func (c *ConsoleLogger) Debug(kvs ...*KeyValue) {
	m := make(map[string]any)
	for _, item := range kvs {
		m[item.Key] = item.Value
	}

	l := template.Bind(c.binder, m)
	log.Println("DEBUG", l)
}

func (c *ConsoleLogger) Info(kvs ...*KeyValue) {
	m := make(map[string]any)
	for _, item := range kvs {
		m[item.Key] = item.Value
	}

	l := template.Bind(c.binder, m)
	log.Println("INFO", l)
}

func (c *ConsoleLogger) Warning(kvs ...*KeyValue) {
	m := make(map[string]any)
	for _, item := range kvs {
		m[item.Key] = item.Value
	}

	l := template.Bind(c.binder, m)
	log.Println("WARNING", l)
}

func (c *ConsoleLogger) Error(err error) {
	log.Println("ERROR", err)
}

func (c *ConsoleLogger) Flush() {

}

func (c *ConsoleTracer) Add(kv *KeyValue) {
	c.data[kv.Key] = kv.Value
}

func (c *ConsoleTracer) Notify() {
	log.Println("TRACE", template.Bind(c.binder, c.data))
}

func (c *ConsoleTracer) NotifyOne(key string) {
	log.Println("TRACE", key, c.data[key])
}

func (c *ConsoleTracer) Flush() {}

func (c *ConsoleTracer) Reset() {
	for key := range c.data {
		delete(c.data, key)
	}
}

func (c *ConsoleMeter) Add(kv *KeyValue) {
	c.data[kv.Key] = kv.Value
}

func (c *ConsoleMeter) Notify() {
	for key, value := range c.data {
		log.Println("TRACE", key, value)
	}
}

func (c *ConsoleMeter) NotifyOne(key string) {
	log.Println("TRACE", key, c.data[key])
}

func (c *ConsoleMeter) Flush() {}

func (c *ConsoleMeter) Reset() {
	for key := range c.data {
		delete(c.data, key)
	}
}
