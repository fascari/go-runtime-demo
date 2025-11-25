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

## Documentation

📚 **API Guide**: [docs/GUIDE.md](./docs/GUIDE.md) - Complete API documentation with scheduler behavior explanations

📖 **OpenAPI Spec**: [docs/openapi.yaml](./docs/openapi.yaml) - OpenAPI 3.0 specification for all endpoints

## Example Usage

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

For an in-depth explanation of Go's scheduler, read the full article on Medium:
- **[Understanding Go's Scheduler: How Goroutine Management Works](https://medium.com/@felipe.ascari_49171/understanding-gos-scheduler-how-goroutine-management-works-65131986ee2c)**

Topics covered: G-M-P model, work-stealing algorithm, goroutine states, preemption, and practical code examples.

---

**Note:** This is an educational project. The blockchain implementation is intentionally simple to focus on demonstrating scheduler behavior.

