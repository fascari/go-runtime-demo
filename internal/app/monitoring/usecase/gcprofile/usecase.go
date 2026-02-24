package gcprofile

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"time"
)

type (
	UseCase struct{}

	ProfileType string

	Input struct {
		DurationSeconds int         `json:"duration_seconds"`
		ProfileType     ProfileType `json:"profile_type"`
	}

	Result struct {
		ProfilePath string  `json:"profile_path"`
		ViewCommand string  `json:"view_command"`
		FileSizeKB  int64   `json:"file_size_kb"`
		DurationMs  float64 `json:"duration_ms"`
		GCRuns      uint32  `json:"gc_runs"`
		HeapAllocMB float64 `json:"heap_alloc_mb"`
		Error       string  `json:"error,omitempty"`
	}
)

const (
	ProfileTypeHeap      ProfileType = "heap"
	ProfileTypeCPU       ProfileType = "cpu"
	ProfileTypeGoroutine ProfileType = "goroutine"
	ProfileTypeAllocs    ProfileType = "allocs"
)

func New() UseCase {
	return UseCase{}
}

func (uc UseCase) Execute(_ context.Context, input Input) Result {
	if input.DurationSeconds <= 0 {
		input.DurationSeconds = 5
	}
	if input.ProfileType == "" {
		input.ProfileType = ProfileTypeHeap
	}

	start := time.Now()
	var memBefore runtime.MemStats
	runtime.ReadMemStats(&memBefore)

	// Create temp directory for profiles
	tmpDir := os.TempDir()
	var profilePath string
	var err error

	switch input.ProfileType {
	case ProfileTypeHeap:
		profilePath, err = uc.createHeapProfile(tmpDir)
	case ProfileTypeCPU:
		profilePath, err = uc.createCPUProfile(tmpDir, input.DurationSeconds)
	case ProfileTypeGoroutine:
		profilePath, err = uc.createGoroutineProfile(tmpDir)
	case ProfileTypeAllocs:
		profilePath, err = uc.createAllocsProfile(tmpDir)
	default:
		profilePath, err = uc.createHeapProfile(tmpDir)
	}

	duration := time.Since(start)

	var memAfter runtime.MemStats
	runtime.ReadMemStats(&memAfter)

	if err != nil {
		return Result{
			Error: err.Error(),
		}
	}

	// Get file size
	var fileSize int64
	if info, err := os.Stat(profilePath); err == nil {
		fileSize = info.Size() / 1024 // KB
	}

	return Result{
		ProfilePath: profilePath,
		ViewCommand: "go tool pprof -http=:8080 " + profilePath,
		FileSizeKB:  fileSize,
		DurationMs:  float64(duration.Microseconds()) / 1000,
		GCRuns:      memAfter.NumGC - memBefore.NumGC,
		HeapAllocMB: float64(memAfter.Alloc) / 1024 / 1024,
	}
}

func (uc UseCase) createHeapProfile(tmpDir string) (string, error) {
	// Force GC before taking heap profile
	runtime.GC()

	path := filepath.Join(tmpDir, "heap.pprof")
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if err := pprof.WriteHeapProfile(f); err != nil {
		return "", err
	}

	return path, nil
}

func (uc UseCase) createCPUProfile(tmpDir string, durationSeconds int) (string, error) {
	path := filepath.Join(tmpDir, "cpu.pprof")
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if err := pprof.StartCPUProfile(f); err != nil {
		return "", err
	}

	time.Sleep(time.Duration(durationSeconds) * time.Second)
	pprof.StopCPUProfile()

	return path, nil
}

func (uc UseCase) createGoroutineProfile(tmpDir string) (string, error) {
	path := filepath.Join(tmpDir, "goroutine.pprof")
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	p := pprof.Lookup("goroutine")
	if p == nil {
		return "", os.ErrNotExist
	}

	if err := p.WriteTo(f, 0); err != nil {
		return "", err
	}

	return path, nil
}

func (uc UseCase) createAllocsProfile(tmpDir string) (string, error) {
	path := filepath.Join(tmpDir, "allocs.pprof")
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	p := pprof.Lookup("allocs")
	if p == nil {
		return "", os.ErrNotExist
	}

	if err := p.WriteTo(f, 0); err != nil {
		return "", err
	}

	return path, nil
}
