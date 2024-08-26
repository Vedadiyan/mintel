package mintel

import "time"

type (
	Level string
	Type  string

	Metadata map[string]string

	CreateFunc func(metadata Metadata) Telemetry

	KeyValue struct {
		Key   string
		Value any
	}

	Writer interface {
		Add(...*KeyValue) Writer
		Flush()
	}

	Telemetry interface {
		Logger() Writer
		Tracer() Writer
		Meter() Writer
		Close()
	}

	TelemetryOpt func(Telemetry)
)

var (
	_clients map[string]CreateFunc
)

func init() {
	_clients = make(map[string]CreateFunc)
}

func TraceRef[T any](name string, ref *T) TelemetryOpt {
	return func(t Telemetry) {
		t.Tracer().Add(KV(name, ref))
	}
}

func Trace[T any](name string, v T) TelemetryOpt {
	return func(t Telemetry) {
		t.Tracer().Add(KV(name, v))
	}
}

func MeasureRef[T int | int16 | int32 | int64 | int8 | byte | uint | uint16 | uint32 | uint64 | float32 | float64](name string, ref *T) TelemetryOpt {
	return func(t Telemetry) {
		t.Meter().Add(KV(name, ref))
	}
}

func Measure[T int | int16 | int32 | int64 | int8 | byte | uint | uint16 | uint32 | uint64 | float32 | float64](name string, v T) TelemetryOpt {
	return func(t Telemetry) {
		t.Meter().Add(KV(name, v))
	}
}

func KV(key string, value any) *KeyValue {
	kv := new(KeyValue)
	kv.Key = key
	kv.Value = value
	return kv
}

func Open(name string, metadata Metadata, opts ...TelemetryOpt) Telemetry {
	var c Telemetry
	if v, ok := _clients[name]; ok {
		c = v(metadata)
		for _, opt := range opts {
			opt(c)
		}
		return c
	}
	return nil
}

func Register(name string, fn CreateFunc) {
	_clients[name] = fn
}

func Verbose() *KeyValue {
	return KV("LEVEL", "Verbose")
}

func Info() *KeyValue {
	return KV("LEVEL", "INFO")
}

func Debug() *KeyValue {
	return KV("LEVEL", "DEBUG")
}

func Warn() *KeyValue {
	return KV("LEVEL", "WARN")
}

func Error() *KeyValue {
	return KV("LEVEL", "ERROR")
}

func Timestamp() *KeyValue {
	return KV("TIMESTAMP", time.Now().UnixNano())
}

func Begin() *KeyValue {
	return KV("PHASE", "BEGIN")
}

func InProgress() *KeyValue {
	return KV("PHASE", "IN_PROGRESS")
}

func End() *KeyValue {
	return KV("PHASE", "END")
}
