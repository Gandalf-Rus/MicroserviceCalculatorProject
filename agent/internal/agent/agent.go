package agent

import (
	"MicroserviceCalculatorProject/orchestrator/pkg/collection"
	"encoding/json"
	"log"
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
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	forever := make(chan bool)
	agent.makeConsumer(conn)

	<-forever
}

func (agent Agent) makeConsumer(conn *amqp.Connection) {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

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
			log.Printf("Received a message: %s", d.Body)

			var message collection.AgentsTask
			err := json.Unmarshal(d.Body, &message)
			failOnError(err, "Failed to unmarshal JSON")

			go func() {
				result := calculateSubexpression(message)
				defer func() {
					agent.FreeWorkers++
					log.Printf("%v, %v", message.ExpressionID, result)
				}()
			}()

			// Подтверждаем получение сообщения
			err = d.Ack(false)
			if err != nil {
				log.Fatalf("Failed to acknowledge message: %v", err)
			}
		}
	}

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func calculateSubexpression(data collection.AgentsTask) float64 {
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

	return result
}
