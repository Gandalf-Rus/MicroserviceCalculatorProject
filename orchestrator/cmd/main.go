package main

import "MicroserviceCalculatorProject/orchestrator/internal/server"

func main() {
	mainOrchestrator := server.New()
	mainOrchestrator.Run()
}
