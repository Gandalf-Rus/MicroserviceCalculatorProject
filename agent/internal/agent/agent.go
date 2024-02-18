package agent

import (
	"MicroserviceCalculatorProject/agent/pkg"
	"MicroserviceCalculatorProject/orchestrator/pkg/collection"
	"MicroserviceCalculatorProject/orchestrator/pkg/logger"
	"encoding/json"
	"log"
	"math"
	"time"

	"github.com/streadway/amqp"
)

type Agent struct {
	CountOfWorkers int
	FreeWorkers    int
}

func New(workers int) *Agent {
	return &Agent{
		CountOfWorkers: workers,
		FreeWorkers:    workers,
	}
}

func (agent *Agent) Run() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	logger.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	logger.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	subexpressionsChan := make(chan pkg.Response)

	go agent.makeConsumer(ch, subexpressionsChan)

	for {
		log.Printf("Im in loop\n")
		message := <-subexpressionsChan
		log.Printf("%v\n", message)
		sendResponse(ch, message)
	}
}

func (agent Agent) makeConsumer(ch *amqp.Channel, resultsChan chan pkg.Response) {
	queue, err := ch.QueueDeclare(
		"UnresolvedTasks", // Queue name
		true,              // Durable
		false,             // Delete when unused
		false,             // Exclusive
		false,             // No-wait
		nil,               // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		queue.Name, // Queue
		"",         // Consumer
		false,      // Auto-ack
		false,      // Exclusive
		false,      // No-local
		false,      // No-wait
		nil,        // Args
	)

	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	for {
		for agent.FreeWorkers > 0 {
			agent.FreeWorkers--
			d := <-msgs

			var message collection.AgentsTask
			err := json.Unmarshal(d.Body, &message)
			logger.FailOnError(err, "Failed to unmarshal JSON")

			go func(task collection.AgentsTask, resultsChan chan pkg.Response) {
				calculateSubexpression(task, resultsChan)
				defer func() {
					agent.FreeWorkers++
				}()
			}(message, resultsChan)

			// Подтверждаем получение сообщения
			err = d.Ack(false)
			if err != nil {
				log.Fatalf("Failed to acknowledge message: %v", err)
			}
		}
	}
}

func sendResponse(ch *amqp.Channel, message pkg.Response) {
	queue, err := ch.QueueDeclare(
		"TasksResults", // Queue name
		true,           // Durable
		false,          // Delete when unused
		false,          // Exclusive
		false,          // No-wait
		nil,            // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}
	body, err := json.Marshal(message)
	if err != nil {
		logger.FailOnError(err, "Failed to marshal JSON")
	}

	err = ch.Publish(
		"",         // Exchange
		queue.Name, // Routing key
		false,      // Mandatory
		false,      // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		logger.FailOnError(err, "Failed to publish")
	}
}

func calculateSubexpression(data collection.AgentsTask, resultChan chan pkg.Response) {
	time.Sleep(time.Second * time.Duration(data.OperatorDuration))
	var result float64
	switch data.Operator {
	case "+":
		result = data.LeftOperand + data.RightOperand
	case "-":
		result = data.LeftOperand - data.RightOperand
	case "*":
		result = data.LeftOperand * data.RightOperand
	case "/":
		result = data.LeftOperand / data.RightOperand
	}

	status := 1
	if result == float64(math.Inf(1)) || result == float64(math.Inf(-1)) {
		result = 0
		status = 3
	}

	response := pkg.Response{
		ExpressionID:        data.ExpressionID,
		SubexpressionNumber: data.SubexpressionNumber,
		Result:              result,
		StatusID:            status,
	}

	log.Printf("%v...\n", response)
	resultChan <- response
}
