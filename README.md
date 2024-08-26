# Mintel - Minimal Telemetry Library

Mintel is a lightweight, flexible telemetry library for Go applications. It provides a standardized interface for logging, tracing, and metrics collection, allowing easy integration with various backend systems.

## Features

- Unified interface for logging, tracing, and metrics
- Support for multiple backend implementations
- Customizable metadata and key-value pairs
- Type-safe measurement and tracing functions
- Easy-to-use fluent API

## Design Philosophy and Advantages

Mintel is designed with a unique philosophy that sets it apart from other telemetry libraries. Here are the key aspects of Mintel's design and their advantages:

### Minimal Core, Powerful Extensibility

Mintel provides a lean, focused core library that handles the essential telemetry operations. This minimalist approach offers several benefits:

1. **Simplicity**: The core API is easy to understand and use, reducing the learning curve for new users.
2. **Flexibility**: The simple core can be easily extended to fit a wide range of use cases.
3. **Low Overhead**: The minimal core ensures that Mintel adds very little overhead to your application.

### Backend-Driven Implementation

One of Mintel's key design decisions is to delegate batching, compatibility, and other advanced features to the backend implementations. This approach offers significant advantages:

1. **Separation of Concerns**: The core library focuses on providing a clean interface, while backends handle the specifics of data processing and transmission.
2. **Optimized Performance**: Backend implementers can optimize batching and other performance aspects based on their specific systems and requirements.
3. **Adaptability**: As new batching techniques or compatibility requirements emerge, backends can be updated without changes to the core library.
4. **Tailored Solutions**: Users can choose or create backends that precisely match their needs, rather than being constrained by a one-size-fits-all approach.

### Unified Interface

Mintel provides a single, unified interface for logging, tracing, and metrics:

1. **Consistency**: Using the same patterns for all telemetry types leads to more consistent and maintainable code.
2. **Simplicity**: Developers only need to learn one interface for all their telemetry needs.

### Type-Safe Measurements

Mintel uses Go's type system to ensure type safety in measurements:

1. **Compile-Time Checks**: Many errors can be caught at compile-time, improving code reliability.
2. **Self-Documenting Code**: The types used in measurements clearly indicate the kind of data expected.

### Comparison with Other Libraries

While more comprehensive libraries like OpenTelemetry offer a wide range of features out-of-the-box, Mintel takes a different approach:

1. **Lean Core**: Mintel's core is much smaller and focused compared to larger libraries.
2. **Extensibility**: While it doesn't provide every feature by default, Mintel's design allows for powerful extensibility through custom backends.
3. **Fine-Grained Control**: Mintel's approach gives developers more control over the specifics of their telemetry implementation.
4. **Lower Complexity**: For projects that don't need all the features of larger libraries, Mintel offers a simpler alternative without sacrificing potential functionality.

 Mintel is designed for developers who value simplicity, flexibility, and control in their telemetry solution. Its unique approach of providing a minimal core with backend-driven implementation offers a compelling alternative to more comprehensive, but also more complex, telemetry libraries.

## Installation

To install Mintel, use `go get`:

```
go get github.com/vedadiyan/mintel
```

## Quick Start

Here's a simple example of how to use Mintel:

```go
package main

import (
    "github.com/vedadiyan/mintel"
)

func main() {
    // Open a new telemetry client with options
    t := mintel.Open("console", mintel.Metadata{
        "service": "my-service",
        "version": "1.0.0",
    },
    mintel.Trace("init", "Starting service"),
    mintel.Measure("startup_time", 1.5))
    defer t.Close()

    // Log a message
    t.Logger().Add(
        mintel.KV("level", "info"),
        mintel.KV("message", "Hello, Mintel!"),
        mintel.KV("timestamp", time.Now()),
    ).Flush()

    // Record a trace
    t.Tracer().Add(
        mintel.KV("operation", "database-query"),
        mintel.KV("query", "SELECT * FROM users"),
    ).Flush()

    // Record a metric
    t.Meter().Add(
        mintel.KV("metric", "active_users"),
        mintel.KV("count", 100),
    ).Flush()
}
```

## Core Concepts

### Telemetry

The `Telemetry` interface is the core of the Mintel library. It provides three main components:

- `Logger()`: For logging events and messages
- `Tracer()`: For tracing operations and recording spans
- `Meter()`: For recording metrics and measurements

### Writer

The `Writer` interface is used by all three components (Logger, Tracer, Meter) to add key-value pairs and flush the data to the backend.

### Metadata

Metadata is a map of string key-value pairs that can be attached to a telemetry client. This is useful for adding context to all telemetry data, such as service name, version, or environment.

### KeyValue

KeyValue structs are used to add data to telemetry records. They consist of a string key and an interface{} value.

## Usage

### Opening a Telemetry Client

To start using Mintel, you need to open a telemetry client:

```go
t := mintel.Open("console", mintel.Metadata{
    "service": "my-service",
    "version": "1.0.0",
},
mintel.Trace("init", "Starting service"),
mintel.Measure("startup_time", 1.5))
defer t.Close()
```

The first argument is the name of the backend to use (e.g., "console", "loki"). The second argument is metadata that will be attached to all telemetry data. Additional arguments are `TelemetryOpt` functions that can be used to add initial traces or measurements.

### Logging

To log a message:

```go
t.Logger().Add(
    mintel.KV("level", "info"),
    mintel.KV("message", "User logged in"),
    mintel.KV("user_id", 12345),
    mintel.KV("timestamp", time.Now()),
).Flush()
```

### Tracing

To record a trace:

```go
t.Tracer().Add(
    mintel.KV("operation", "http-request"),
    mintel.KV("method", "GET"),
    mintel.KV("path", "/api/users"),
    mintel.KV("phase", "begin"),
).Flush()

// ... perform operation ...

t.Tracer().Add(
    mintel.KV("operation", "http-request"),
    mintel.KV("status", 200),
    mintel.KV("phase", "end"),
).Flush()
```

### Metrics

To record a metric:

```go
t.Meter().Add(
    mintel.KV("metric", "request_duration"),
    mintel.KV("value", 123.45),
    mintel.KV("unit", "ms"),
).Flush()
```

## Utility Functions

Mintel provides several utility functions to make it easier to add common key-value pairs:

- `KV(key string, value any)`: Create a KeyValue pair
- `Verbose()`, `Info()`, `Debug()`, `Warn()`, `Error()`: Create KeyValue pairs for log levels
- `Timestamp()`: Create a KeyValue pair with the current timestamp
- `Begin()`, `Exec()`, `End()`: Create KeyValue pairs for operation phases

Mintel also provides TelemetryOpt functions for use with `Open()`:

- `TraceRef(name string, ref *T)`, `Trace(name string, v T)`: Add trace information
- `MeasureRef(name string, ref *T)`, `Measure(name string, v T)`: Add measurements

## Template Writer

Mintel includes a template writer utility that allows you to define log templates. For example:

```go
template := `{
    "name": $(name),
    "age": $(age),
    "address": {
        "street": $(address.street),
        "city": $(address.city)
    }
}`

binder, err := template.Parse(template)
if err != nil {
    // handle error
}

data := map[string]interface{}{
    "name": "John Doe",
    "age": 30,
    "address": map[string]string{
        "street": "123 Main St",
        "city": "Anytown",
    },
}

result := template.Bind(binder, data)
fmt.Println(result)
```

This will output:

```json
{
    "name": "John Doe",
    "age": 30,
    "address": {
        "street": "123 Main St",
        "city": "Anytown"
    }
}
```

## JSON Marshaller

Mintel includes a custom JSON marshaller that can handle reference cycles and other complex structures. To use it:

```go
import "github.com/vedadiyan/mintel/util/json"

data := // ... your data structure ...
jsonBytes := json.Marshal(data)
```

This marshaller will not yield any errors, making it safer to use in logging scenarios where you don't want to risk losing data due to marshalling errors.

## Extending Mintel

You can extend Mintel by implementing your own backend. To do this, create a new type that implements the `Telemetry` interface and register it using the `Register` function:

```go
func init() {
    mintel.Register("my-backend", func(metadata mintel.Metadata) mintel.Telemetry {
        return &MyBackend{metadata: metadata}
    })
}
```

## Contributing

Contributions to Mintel are welcome! Please feel free to submit a Pull Request.

## License

Mintel is released under the MIT License. See the LICENSE file for details.