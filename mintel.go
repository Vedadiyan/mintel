package mintel

type (
	Level string
	Type  string

	Metadata map[string]string

	CreateClient func(Metadata) Telemetry

	KV struct {
		Key   string
		Value any
	}

	Logger interface {
		Debug(...*KV)
		Info(...*KV)
		Warning(...*KV)
		Error(error)
		Flush()
	}

	Tracer interface {
		Add(*KV)
		Notify()
		NotifyOne(string)
		Flush()
		Reset()
	}

	Meter interface {
		Add(*KV)
		Notify()
		NotifyOne(string)
		Flush()
		Reset()
	}

	Telemetry interface {
		Logger() Logger
		Tracer() Tracer
		Meter() Meter
		Close(...any)
	}

	TelemetryOpt func(Telemetry)
)

var (
	_clients map[string]CreateClient
)

func TraceRef[T any](name string, ref *T) TelemetryOpt {
	return func(t Telemetry) {
		t.Tracer().Add(KeyValue(name, ref))
	}
}

func Trace[T any](name string, v T) TelemetryOpt {
	return func(t Telemetry) {
		t.Tracer().Add(KeyValue(name, v))
	}
}

func MeasureRef[T int | int16 | int32 | int64 | int8 | byte | uint | uint16 | uint32 | uint64 | float32 | float64](name string, ref *T) TelemetryOpt {
	return func(t Telemetry) {
		t.Meter().Add(KeyValue(name, ref))
	}
}

func Measure[T int | int16 | int32 | int64 | int8 | byte | uint | uint16 | uint32 | uint64 | float32 | float64](name string, v T) TelemetryOpt {
	return func(t Telemetry) {
		t.Meter().Add(KeyValue(name, v))
	}
}

func KeyValue(key string, value any) *KV {
	kv := new(KV)
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

func RegisterClient(name string, fn CreateClient) {
	_clients[name] = fn
}
