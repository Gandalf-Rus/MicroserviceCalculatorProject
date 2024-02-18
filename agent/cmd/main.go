package main

import (
	"MicroserviceCalculatorProject/agent/internal/agent"
)

func main() {
	agent := agent.New(5)
	agent.Run()
}
