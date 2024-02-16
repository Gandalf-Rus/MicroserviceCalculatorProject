package main

import "MicroserviceCalculatorProject/orchestrator/internal/server"

func main() {
	orchestrator := server.New()
	orchestrator.Run()
}
