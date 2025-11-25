package domain

import (
	"runtime"
	"runtime/debug"
	"time"
)

type (
	RuntimeStats struct {
		NumGoroutine   int           `json:"num_goroutine"`
		NumCPU         int           `json:"num_cpu"`
		MemoryStats    MemoryStats   `json:"memory_stats"`
		GCStats        GCStats       `json:"gc_stats"`
		SchedulerStats SchedulerInfo `json:"scheduler_stats"`
		Timestamp      time.Time     `json:"timestamp"`
	}

	MemoryStats struct {
		AllocMB        uint64 `json:"alloc_mb"`
		TotalAllocMB   uint64 `json:"total_alloc_mb"`
		SysMB          uint64 `json:"sys_mb"`
		NumGC          uint32 `json:"num_gc"`
		HeapAllocMB    uint64 `json:"heap_alloc_mb"`
		HeapObjects    uint64 `json:"heap_objects"`
		HeapIdleMB     uint64 `json:"heap_idle_mb"`
		HeapInUseMB    uint64 `json:"heap_in_use_mb"`
		HeapReleasedMB uint64 `json:"heap_released_mb"`
	}

	GCStats struct {
		LastGC        string  `json:"last_gc"`
		NextGCMB      uint64  `json:"next_gc_mb"`
		PauseTotal    string  `json:"pause_total"`
		NumForcedGC   uint32  `json:"num_forced_gc"`
		GCCPUFraction float64 `json:"gc_cpu_fraction"`
	}

	SchedulerInfo struct {
		NumProcs      int   `json:"num_procs"`
		NumGoroutines int   `json:"num_goroutines"`
		NumCgoCall    int64 `json:"num_cgo_call"`
	}

	GCResult struct {
		GCRunsBefore  uint32  `json:"gc_runs_before"`
		GCRunsAfter   uint32  `json:"gc_runs_after"`
		MemoryFreedMB float64 `json:"memory_freed_mb"`
		AllocBeforeMB float64 `json:"alloc_before_mb"`
		AllocAfterMB  float64 `json:"alloc_after_mb"`
	}

	Monitor struct{}
)

func NewMonitor() *Monitor {
	return &Monitor{}
}

func (m *Monitor) RuntimeStats() RuntimeStats {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	var gcStats debug.GCStats
	debug.ReadGCStats(&gcStats)

	return RuntimeStats{
		NumGoroutine: runtime.NumGoroutine(),
		NumCPU:       runtime.NumCPU(),
		MemoryStats: MemoryStats{
			AllocMB:        memStats.Alloc / 1024 / 1024,
			TotalAllocMB:   memStats.TotalAlloc / 1024 / 1024,
			SysMB:          memStats.Sys / 1024 / 1024,
			NumGC:          memStats.NumGC,
			HeapAllocMB:    memStats.HeapAlloc / 1024 / 1024,
			HeapObjects:    memStats.HeapObjects,
			HeapIdleMB:     memStats.HeapIdle / 1024 / 1024,
			HeapInUseMB:    memStats.HeapInuse / 1024 / 1024,
			HeapReleasedMB: memStats.HeapReleased / 1024 / 1024,
		},
		GCStats: GCStats{
			LastGC:        time.Unix(0, int64(memStats.LastGC)).Format(time.RFC3339),
			NextGCMB:      memStats.NextGC / 1024 / 1024,
			PauseTotal:    time.Duration(memStats.PauseTotalNs).String(),
			NumForcedGC:   memStats.NumForcedGC,
			GCCPUFraction: memStats.GCCPUFraction,
		},
		SchedulerStats: SchedulerInfo{
			NumProcs:      runtime.GOMAXPROCS(0),
			NumGoroutines: runtime.NumGoroutine(),
			NumCgoCall:    runtime.NumCgoCall(),
		},
		Timestamp: time.Now(),
	}
}

func (m *Monitor) ForceGC() GCResult {
	var before, after runtime.MemStats

	runtime.ReadMemStats(&before)
	runtime.GC()
	runtime.ReadMemStats(&after)

	return GCResult{
		GCRunsBefore:  before.NumGC,
		GCRunsAfter:   after.NumGC,
		MemoryFreedMB: float64(before.Alloc-after.Alloc) / 1024 / 1024,
		AllocBeforeMB: float64(before.Alloc) / 1024 / 1024,
		AllocAfterMB:  float64(after.Alloc) / 1024 / 1024,
	}
}
