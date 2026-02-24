package main

import (
	"fmt"
	"log"
	"runtime"

	addblockhandler "go-runtime-demo/internal/app/blockchain/handler/addblock"
	listblockshandler "go-runtime-demo/internal/app/blockchain/handler/listblocks"
	mineparallelhandler "go-runtime-demo/internal/app/blockchain/handler/mineparallel"
	stresstesthandler "go-runtime-demo/internal/app/blockchain/handler/stresstest"
	gcbenchmarkhandler "go-runtime-demo/internal/app/monitoring/handler/gcbenchmark"
	gcfinalizershandler "go-runtime-demo/internal/app/monitoring/handler/gcfinalizers"
	gcmetricshandler "go-runtime-demo/internal/app/monitoring/handler/gcmetrics"
	gcprofilehandler "go-runtime-demo/internal/app/monitoring/handler/gcprofile"
	statshandler "go-runtime-demo/internal/app/monitoring/handler/stats"

	blockchaindomain "go-runtime-demo/internal/app/blockchain/domain"
	addblockusecase "go-runtime-demo/internal/app/blockchain/usecase/addblock"
	listblocksusecase "go-runtime-demo/internal/app/blockchain/usecase/listblocks"
	mineparallelusecase "go-runtime-demo/internal/app/blockchain/usecase/mineparallel"
	stresstestusecase "go-runtime-demo/internal/app/blockchain/usecase/stresstest"
	gcbenchmarkusecase "go-runtime-demo/internal/app/monitoring/usecase/gcbenchmark"
	gcfinalizersusecase "go-runtime-demo/internal/app/monitoring/usecase/gcfinalizers"
	gcmetricsusecase "go-runtime-demo/internal/app/monitoring/usecase/gcmetrics"
	gcprofileusecase "go-runtime-demo/internal/app/monitoring/usecase/gcprofile"

	monitoringdomain "go-runtime-demo/internal/app/monitoring/domain"
	statsusecase "go-runtime-demo/internal/app/monitoring/usecase/stats"

	httpserver "go-runtime-demo/pkg/http"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	printSchedulerInfo()

	blockchain := blockchaindomain.NewBlockchain(4)
	monitor := monitoringdomain.NewMonitor()

	// Blockchain use cases
	addBlockUC := addblockusecase.New(blockchain)
	listBlocksUC := listblocksusecase.New(blockchain)
	mineParallelUC := mineparallelusecase.New(blockchain)
	stressTestUC := stresstestusecase.New()

	// Monitoring use cases
	statsUC := statsusecase.New(monitor)
	gcBenchmarkUC := gcbenchmarkusecase.New()
	gcFinalizersUC := gcfinalizersusecase.New()
	gcMetricsUC := gcmetricsusecase.New()
	gcProfileUC := gcprofileusecase.New()

	// Handlers
	addBlockHandler := addblockhandler.NewHandler(addBlockUC)
	listBlocksHandler := listblockshandler.NewHandler(listBlocksUC)
	mineParallelHandler := mineparallelhandler.NewHandler(mineParallelUC)
	stressTestHandler := stresstesthandler.NewHandler(stressTestUC)
	statsHandler := statshandler.NewHandler(statsUC)
	gcBenchmarkHandler := gcbenchmarkhandler.NewHandler(gcBenchmarkUC)
	gcFinalizersHandler := gcfinalizershandler.NewHandler(gcFinalizersUC)
	gcMetricsHandler := gcmetricshandler.NewHandler(gcMetricsUC)
	gcProfileHandler := gcprofilehandler.NewHandler(gcProfileUC)

	server := httpserver.NewServer("8080")
	router := server.Router()

	// Blockchain endpoints
	addblockhandler.RegisterEndpoint(router, addBlockHandler)
	listblockshandler.RegisterEndpoint(router, listBlocksHandler)
	mineparallelhandler.RegisterEndpoint(router, mineParallelHandler)
	stresstesthandler.RegisterEndpoint(router, stressTestHandler)

	// Monitoring endpoints
	statshandler.RegisterEndpoint(router, statsHandler)
	gcbenchmarkhandler.RegisterEndpoint(router, gcBenchmarkHandler)
	gcfinalizershandler.RegisterEndpoint(router, gcFinalizersHandler)
	gcmetricshandler.RegisterEndpoint(router, gcMetricsHandler)
	gcprofilehandler.RegisterEndpoint(router, gcProfileHandler)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

func printSchedulerInfo() {
	fmt.Printf("=== Go Scheduler Configuration ===\n")
	fmt.Printf("NumCPU: %d\n", runtime.NumCPU())
	fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))
	fmt.Printf("NumGoroutine: %d\n", runtime.NumGoroutine())
	fmt.Printf("==================================\n\n")
}
