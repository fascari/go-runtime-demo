package gcmetrics

import (
	"context"
	"runtime/metrics"
)

type (
	UseCase struct{}

	// MetricsResult holds the collected runtime/metrics data
	MetricsResult struct {
		// GC Heap Metrics
		GCHeapAllocsBytes   uint64 `json:"/gc/heap/allocs:bytes"`
		GCHeapAllocsObjects uint64 `json:"/gc/heap/allocs:objects"`
		GCHeapFreesBytes    uint64 `json:"/gc/heap/frees:bytes"`
		GCHeapFreesObjects  uint64 `json:"/gc/heap/frees:objects"`
		GCHeapGoalBytes     uint64 `json:"/gc/heap/goal:bytes"`
		GCHeapObjects       uint64 `json:"/gc/heap/objects:objects"`
		GCHeapTinyAllocs    uint64 `json:"/gc/heap/tiny/allocs:objects"`

		// GC Cycle Metrics
		GCCyclesTotal  uint64 `json:"/gc/cycles/total:gc-cycles"`
		GCCyclesForced uint64 `json:"/gc/cycles/forced-gc:gc-cycles"`

		// GC Pause Metrics
		GCPauseTotalSeconds float64 `json:"/gc/pause/total:seconds"`
		GCPauseLatencyP50   float64 `json:"/gc/pause/p50:seconds"`
		GCPauseLatencyP99   float64 `json:"/gc/pause/p99:seconds"`

		// Memory Metrics
		MemoryClassesTotal uint64 `json:"/memory/classes/total:bytes"`
		MemoryClassesHeap  uint64 `json:"/memory/classes/heap/released:bytes"`

		// Scheduler Metrics
		SchedGoroutines     uint64  `json:"/sched/goroutines:goroutines"`
		SchedLatencySeconds float64 `json:"/sched/latency:seconds"`

		// Computed metrics
		AllocRateMBPerSec float64 `json:"alloc_rate_mb_per_sec"`
		GoVersion         string  `json:"go_version"`
	}
)

func New() UseCase {
	return UseCase{}
}

func (uc UseCase) Execute(_ context.Context) MetricsResult {
	// Define the metrics we want to collect
	sample := []metrics.Sample{
		{Name: "/gc/heap/allocs:bytes"},
		{Name: "/gc/heap/allocs:objects"},
		{Name: "/gc/heap/frees:bytes"},
		{Name: "/gc/heap/frees:objects"},
		{Name: "/gc/heap/goal:bytes"},
		{Name: "/gc/heap/objects:objects"},
		{Name: "/gc/heap/tiny/allocs:objects"},
		{Name: "/gc/cycles/total:gc-cycles"},
		{Name: "/gc/cycles/forced-gc:gc-cycles"},
		{Name: "/gc/pause/total:seconds"},
		{Name: "/memory/classes/total:bytes"},
		{Name: "/memory/classes/heap/released:bytes"},
		{Name: "/sched/goroutines:goroutines"},
	}

	metrics.Read(sample)

	result := MetricsResult{
		GoVersion: runtimeVersion(),
	}

	for i, s := range sample {
		if s.Value.Kind() == metrics.KindBad {
			continue
		}

		switch sample[i].Name {
		case "/gc/heap/allocs:bytes":
			result.GCHeapAllocsBytes = s.Value.Uint64()
		case "/gc/heap/allocs:objects":
			result.GCHeapAllocsObjects = s.Value.Uint64()
		case "/gc/heap/frees:bytes":
			result.GCHeapFreesBytes = s.Value.Uint64()
		case "/gc/heap/frees:objects":
			result.GCHeapFreesObjects = s.Value.Uint64()
		case "/gc/heap/goal:bytes":
			result.GCHeapGoalBytes = s.Value.Uint64()
		case "/gc/heap/objects:objects":
			result.GCHeapObjects = s.Value.Uint64()
		case "/gc/heap/tiny/allocs:objects":
			result.GCHeapTinyAllocs = s.Value.Uint64()
		case "/gc/cycles/total:gc-cycles":
			result.GCCyclesTotal = s.Value.Uint64()
		case "/gc/cycles/forced-gc:gc-cycles":
			result.GCCyclesForced = s.Value.Uint64()
		case "/gc/pause/total:seconds":
			result.GCPauseTotalSeconds = s.Value.Float64()
		case "/memory/classes/total:bytes":
			result.MemoryClassesTotal = s.Value.Uint64()
		case "/memory/classes/heap/released:bytes":
			result.MemoryClassesHeap = s.Value.Uint64()
		case "/sched/goroutines:goroutines":
			result.SchedGoroutines = s.Value.Uint64()
		}
	}

	// Compute allocation rate (approximate)
	if result.GCHeapAllocsBytes > 0 {
		result.AllocRateMBPerSec = float64(result.GCHeapAllocsBytes) / 1024 / 1024
	}

	return result
}

func runtimeVersion() string {
	return "go1.26"
}
