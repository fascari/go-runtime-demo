package main

import (
	"fmt"
	"log"
	"runtime"

	addblockhandler "go-runtime-demo/internal/app/blockchain/handler/addblock"
	listblockshandler "go-runtime-demo/internal/app/blockchain/handler/listblocks"
	mineparallelhandler "go-runtime-demo/internal/app/blockchain/handler/mineparallel"
	stresstesthandler "go-runtime-demo/internal/app/blockchain/handler/stresstest"
	statshandler "go-runtime-demo/internal/app/monitoring/handler/stats"

	blockchaindomain "go-runtime-demo/internal/app/blockchain/domain"
	addblockusecase "go-runtime-demo/internal/app/blockchain/usecase/addblock"
	listblocksusecase "go-runtime-demo/internal/app/blockchain/usecase/listblocks"
	mineparallelusecase "go-runtime-demo/internal/app/blockchain/usecase/mineparallel"
	stresstestusecase "go-runtime-demo/internal/app/blockchain/usecase/stresstest"

	monitoringdomain "go-runtime-demo/internal/app/monitoring/domain"
	statsusecase "go-runtime-demo/internal/app/monitoring/usecase/stats"

	httpserver "go-runtime-demo/pkg/http"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	printSchedulerInfo()

	// Initialize domain objects
	blockchain := blockchaindomain.NewBlockchain(4)
	monitor := monitoringdomain.NewMonitor()

	// Initialize use cases
	addBlockUC := addblockusecase.New(blockchain)
	listBlocksUC := listblocksusecase.New(blockchain)
	mineParallelUC := mineparallelusecase.New(blockchain)
	stressTestUC := stresstestusecase.New()
	statsUC := statsusecase.New(monitor)

	// Initialize handlers
	addBlockHandler := addblockhandler.NewHandler(addBlockUC)
	listBlocksHandler := listblockshandler.NewHandler(listBlocksUC)
	mineParallelHandler := mineparallelhandler.NewHandler(mineParallelUC)
	stressTestHandler := stresstesthandler.NewHandler(stressTestUC)
	statsHandler := statshandler.NewHandler(statsUC)

	// Setup server and register endpoints
	server := httpserver.NewServer("8080")
	router := server.Router()

	addblockhandler.RegisterEndpoint(router, addBlockHandler)
	listblockshandler.RegisterEndpoint(router, listBlocksHandler)
	mineparallelhandler.RegisterEndpoint(router, mineParallelHandler)
	stresstesthandler.RegisterEndpoint(router, stressTestHandler)
	statshandler.RegisterEndpoint(router, statsHandler)

	// Start server
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
