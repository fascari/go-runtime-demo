# Go Runtime Demo

A practical demonstration of Go's runtime internals using a simple blockchain implementation. This project serves as a hands-on exploration of the Go runtime, including the scheduler, garbage collector, memory management, and other runtime concepts.

## Purpose

This project was created to understand how Go's runtime works in real-world scenarios. By implementing CPU-intensive operations (blockchain mining) and monitoring runtime behavior, we can observe:

### Scheduler Behavior
- How goroutines are distributed across logical processors (Ps)
- Work-stealing algorithm in action
- Cooperative scheduling with `runtime.Gosched()`
- Goroutine state transitions
- The impact of concurrent operations on system resources

### Garbage Collection (Go 1.26 "Green Tea")
- GC metrics via `runtime/metrics` API
- Allocation patterns (short-lived, long-lived, mixed)
- GC pause times and CPU fraction
- Finalizer behavior
- Profile generation (heap, CPU, goroutine, allocs)

The **monitoring module** provides real-time visibility into Go's runtime behavior, exposing metrics like:
- Active goroutines count
- GOMAXPROCS settings
- Memory allocation patterns
- GC pause times and cycles
- Scheduler statistics
- Heap allocation/frees
- GC CPU fraction

This visibility helps understand the relationship between code patterns and runtime behavior, making it easier to write efficient concurrent Go programs.

## Quick Start

```bash
# Clone and run
git clone https://github.com/fascari/go-runtime-demo
cd go-runtime-demo
go run ./cmd/api
```

Server starts on `http://localhost:8080`

## Documentation

📚 **API Guide**: [docs/GUIDE.md](./docs/GUIDE.md) - Complete API documentation with scheduler behavior explanations

📖 **OpenAPI Spec**: [docs/openapi.yaml](./docs/openapi.yaml) - OpenAPI 3.0 specification for all endpoints

## API Endpoints

### Blockchain Operations

**Add a block:**
```bash
curl -X POST http://localhost:8080/blocks \
  -H "Content-Type: application/json" \
  -d '{"data":"Transaction data"}'
```

**Mine blocks in parallel:**
```bash
curl -X POST http://localhost:8080/mine \
  -H "Content-Type: application/json" \
  -d '{"data":"Parallel mining","goroutines":4}'
```

**List blocks:**
```bash
curl http://localhost:8080/blocks | jq .
```

### Stress Testing

**Stress test with GC metrics:**
```bash
# Short-lived allocations (most common pattern)
curl -X POST http://localhost:8080/stress \
  -H "Content-Type: application/json" \
  -d '{"allocations": 100, "goroutines": 4, "pattern": "short-lived"}'

# Long-lived allocations (simulates caching)
curl -X POST http://localhost:8080/stress \
  -H "Content-Type: application/json" \
  -d '{"allocations": 100, "goroutines": 4, "pattern": "long-lived"}'

# Mixed allocation pattern
curl -X POST http://localhost:8080/stress \
  -H "Content-Type: application/json" \
  -d '{"allocations": 100, "goroutines": 4, "pattern": "mixed"}'
```

### GC Monitoring

**Get GC metrics (runtime/metrics API):**
```bash
curl http://localhost:8080/gc/metrics | jq .
```

**Controlled GC benchmark:**
```bash
curl -X POST http://localhost:8080/gc/benchmark \
  -H "Content-Type: application/json" \
  -d '{"allocations": 10000, "size_kb": 1, "pattern": "short-lived"}'
```

**Test finalizers:**
```bash
curl -X POST http://localhost:8080/gc/finalizers \
  -H "Content-Type: application/json" \
  -d '{"count": 100, "trigger_gc": true}'
```

**Generate profiles:**
```bash
# Heap profile
curl -X POST http://localhost:8080/gc/profile \
  -H "Content-Type: application/json" \
  -d '{"profile_type": "heap"}'

# CPU profile (runs for specified duration)
curl -X POST http://localhost:8080/gc/profile \
  -H "Content-Type: application/json" \
  -d '{"profile_type": "cpu", "duration_seconds": 5}'

# View profile
go tool pprof -http=:8080 /tmp/heap.pprof
```

**View scheduler statistics:**
```bash
curl http://localhost:8080/stats | jq .
```

### Comparing GC Settings

```bash
# Default (GOGC=100)
GOGC=100 go run ./cmd/api &

# More aggressive (lower memory, more CPU)
GOGC=50 go run ./cmd/api &

# Less aggressive (higher memory, less CPU)
GOGC=200 go run ./cmd/api &
```

## Related Articles

- **[Understanding Go's Scheduler: How Goroutine Management Works](https://medium.com/@felipe.ascari_49171/understanding-gos-scheduler-how-goroutine-management-works-65131986ee2c)** - G-M-P model, work-stealing algorithm, goroutine states, preemption

- **[Green Tea: Understanding Go's Garbage Collector](https://medium.com/@felipe.ascari_49171/green-tea-understanding-gos-garbage-collector-21cc1bc08725)** - GC evolution, runtime/metrics, allocation patterns, profiling

Topics covered: G-M-P model, work-stealing algorithm, goroutine states, preemption, GC internals, and practical code examples.

---

**Note:** This is an educational project. The blockchain implementation is intentionally simple to focus on demonstrating runtime behavior.

