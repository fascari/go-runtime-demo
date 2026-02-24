package gcfinalizers

import (
	"context"
	"runtime"
	"sync/atomic"
	"time"
)

type (
	UseCase struct{}

	Input struct {
		Count     int  `json:"count"`
		TriggerGC bool `json:"trigger_gc"`
	}

	Result struct {
		ObjectsCreated     int     `json:"objects_created"`
		FinalizersExecuted int32   `json:"finalizers_executed"`
		DurationMs         float64 `json:"duration_ms"`
		GCRuns             uint32  `json:"gc_runs"`
		Warning            string  `json:"warning,omitempty"`
		FinalizerLatencyMs float64 `json:"finalizer_latency_ms"`
	}
)

// objectWithFinalizer is a type that will have a finalizer attached
type objectWithFinalizer struct {
	id int
}

func New() UseCase {
	return UseCase{}
}

func (uc UseCase) Execute(_ context.Context, input Input) Result {
	if input.Count <= 0 {
		input.Count = 100
	}

	var memBefore, memAfter runtime.MemStats
	runtime.ReadMemStats(&memBefore)

	var finalizerCount atomic.Int32
	start := time.Now()

	// Create objects with finalizers
	for i := 0; i < input.Count; i++ {
		obj := &objectWithFinalizer{id: i}
		runtime.SetFinalizer(obj, func(o *objectWithFinalizer) {
			finalizerCount.Add(1)
		})
		// obj goes out of scope here and becomes eligible for GC
	}

	// Trigger GC if requested
	if input.TriggerGC {
		runtime.GC()
		runtime.GC() // Run twice to ensure finalizers are executed
	}

	// Wait a bit for finalizers to complete
	time.Sleep(10 * time.Millisecond)

	duration := time.Since(start)
	runtime.ReadMemStats(&memAfter)

	executedCount := finalizerCount.Load()

	result := Result{
		ObjectsCreated:     input.Count,
		FinalizersExecuted: executedCount,
		DurationMs:         float64(duration.Microseconds()) / 1000,
		GCRuns:             memAfter.NumGC - memBefore.NumGC,
	}

	// Calculate finalizer latency estimate
	if executedCount > 0 {
		result.FinalizerLatencyMs = result.DurationMs / float64(executedCount)
	}

	// Add warning if not all finalizers ran
	if int(executedCount) < input.Count {
		result.Warning = "Finalizers are not guaranteed to run immediately; some may execute on future GC cycles"
	}

	return result
}
