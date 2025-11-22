# Go Runtime Demo

A practical demonstration of Go's runtime internals using a simple blockchain implementation. This project serves as a hands-on exploration of the Go runtime, including the scheduler, garbage collector, memory management, and other runtime concepts.

## Purpose

This project was created to understand how Go's scheduler works in real-world scenarios. By implementing CPU-intensive operations (blockchain mining) and monitoring runtime behavior, we can observe:

- How goroutines are distributed across logical processors (Ps)
- Work-stealing algorithm in action
- Cooperative scheduling with `runtime.Gosched()`
- Goroutine state transitions
- The impact of concurrent operations on system resources

The **monitoring module** provides real-time visibility into Go's scheduler behavior, exposing metrics like:
- Active goroutines count
- GOMAXPROCS settings
- Memory allocation patterns
- GC pause times
- Scheduler statistics

This visibility helps understand the relationship between code patterns and scheduler behavior, making it easier to write efficient concurrent Go programs.

## Quick Start

```bash
# Clone and run
git clone https://github.com/fascari/go-runtime-demo
cd go-runtime-demo
go run ./cmd/api
```

Server starts on `http://localhost:8080`

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/blocks` | List all blocks in the chain |
| POST | `/blocks` | Add a single block |
| POST | `/mine` | Mine blocks using multiple goroutines |
| POST | `/stress` | Run stress test to observe scheduler behavior |
| GET | `/stats` | Get real-time runtime and scheduler statistics |

### Examples

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

**View scheduler statistics:**
```bash
curl http://localhost:8080/stats | jq .
```

## Related Article

For an in-depth explanation of Go's scheduler, check out the accompanying article:
- [Understanding Go's Scheduler](./ARTICLE_SCHEDULER_EN.md)

Topics covered: G-M-P model, work-stealing algorithm, goroutine states, preemption, and practical code examples.

---

**Note:** This is an educational project. The blockchain implementation is intentionally simple to focus on demonstrating scheduler behavior.

