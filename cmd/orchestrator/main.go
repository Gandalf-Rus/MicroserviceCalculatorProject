package main

import "MicroserviceCalculatorProject/internal/orchestrator"

func main() {
	orch := orchestrator.New()
	orch.Run()
}
