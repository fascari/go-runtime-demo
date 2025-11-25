# API Guide

## Overview

This API demonstrates the Go runtime scheduler behavior through a blockchain implementation. It provides endpoints to interact with a blockchain and monitor Go runtime metrics such as goroutines, memory usage, and garbage collection.

For detailed endpoint specifications, request/response schemas, and examples, see the [OpenAPI specification](./openapi.yaml).

## Running the Server

```bash
go run cmd/api/main.go
```

The server starts on port 8080 and displays initial scheduler configuration:

```
=== Go Scheduler Configuration ===
NumCPU: 8
GOMAXPROCS: 8
NumGoroutine: 1
==================================

Server starting on :8080
```

## Available Endpoints

- `GET /stats` - Get runtime statistics
- `POST /blocks` - Add a block to the blockchain
- `GET /blocks` - List all blocks
- `POST /mine` - Mine blocks in parallel
- `POST /stress` - Run stress test

## Understanding Go Scheduler Metrics

### Goroutines

**num_goroutine**: Number of currently running goroutines. This increases when:
- New goroutines are spawned (mining, stress tests)
- Background operations are running
- HTTP handlers are processing requests

The value should return to baseline when operations complete. If it stays high, it indicates goroutines that haven't finished or are blocked.

### Processors (P)

**num_procs** (GOMAXPROCS): Number of OS threads that can execute Go code simultaneously. This is typically set to the number of CPU cores. The Go scheduler distributes goroutines across these processors.

**NumCPU**: Total number of logical CPU cores available to the process.

### Memory Metrics

**alloc_mb**: Current allocated memory in MB. This shows memory actively used by the application.

**heap_alloc_mb**: Memory allocated on the heap. Similar to alloc_mb but excludes stack allocations.

**heap_in_use_mb**: Memory currently in use by the heap. This includes allocated objects and unused space in heap spans.

**heap_idle_mb**: Memory in the heap that is not currently being used and could be returned to the OS.

**heap_released_mb**: Memory that has been returned to the OS.

**sys_mb**: Total memory obtained from the OS. This includes heap, stack, and internal Go structures.

**total_alloc_mb**: Cumulative bytes allocated throughout the program's lifetime. This always increases.

### Garbage Collection

**num_gc**: Number of completed GC cycles. This increases when the heap grows beyond the GC target.

**next_gc_mb**: Heap size target for the next GC cycle. When heap_alloc_mb approaches this value, GC will trigger.

**pause_total**: Cumulative time the application was paused for GC stop-the-world phases.

**gc_cpu_fraction**: Fraction of CPU time used by GC. Values close to 0 are ideal. Higher values indicate GC pressure.

**num_forced_gc**: Number of GC cycles triggered manually via runtime.GC().

**last_gc**: Timestamp of the last GC cycle.

## Observing Scheduler Behavior

### Mining a Single Block

When you call POST /blocks, the operation runs synchronously in a single goroutine. Observe:
- num_goroutine increases by 1-2 during mining
- CPU usage on one core
- Memory remains relatively stable

### Parallel Mining

When you call POST /mine with multiple goroutines, observe:
- num_goroutine spikes to the requested number
- CPU usage distributed across multiple cores
- Duration decreases with more goroutines (up to NumCPU)
- Work-stealing: goroutines are distributed across available Ps

### Stress Test

When you call POST /stress, observe:
- num_goroutine increases to requested value
- memory_delta_mb shows allocated memory during test
- num_gc increases if GC is triggered
- gc_cpu_fraction may increase under memory pressure
- num_goroutines returns to baseline after completion

### Memory Pressure

During heavy allocations (high stress test values):
- heap_alloc_mb increases
- When approaching next_gc_mb, GC triggers
- heap_idle_mb decreases as heap grows
- pause_total increases with each GC cycle
- num_gc increments

## Example Workflow

```bash
# Start monitoring stats
curl http://localhost:8080/stats

# Add a single block and observe goroutine behavior
curl -X POST http://localhost:8080/blocks \
  -H "Content-Type: application/json" \
  -d '{"data":"Test block"}'

# Mine with 4 goroutines and compare duration
curl -X POST http://localhost:8080/mine \
  -H "Content-Type: application/json" \
  -d '{"data":"Parallel test","goroutines":4}'

# Run stress test to observe GC behavior
curl -X POST http://localhost:8080/stress \
  -H "Content-Type: application/json" \
  -d '{"allocations":100,"goroutines":10}'

# Check final stats to see GC impact
curl http://localhost:8080/stats
```

## Key Observations

1. **Goroutine Scheduling**: The scheduler automatically distributes goroutines across available Ps. You don't need to manually assign goroutines to threads.

2. **Work Stealing**: When a P finishes its work, it can steal goroutines from other Ps' queues, maintaining load balance.

3. **GC Triggering**: GC runs automatically when heap memory approaches the next_gc_mb target. It's adaptive based on allocation rate.

4. **Memory Allocation**: Go uses a concurrent mark-and-sweep collector. Most GC work happens concurrently with the application.

5. **Goroutine Overhead**: Each goroutine has minimal overhead (~2KB stack initially). You can create thousands efficiently.

6. **CPU Utilization**: With proper goroutine usage, Go automatically utilizes all available CPU cores without manual thread management.

