package main

import "MicroserviceCalculatorProject/internal/agent"

func main() {
	agent := agent.New(5)
	agent.Run()
}
