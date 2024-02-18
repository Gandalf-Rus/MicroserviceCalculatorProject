package server

import (
	"MicroserviceCalculatorProject/agent/pkg"
	"MicroserviceCalculatorProject/orchestrator/internal/database"
	"MicroserviceCalculatorProject/orchestrator/pkg/collection"
	"MicroserviceCalculatorProject/orchestrator/pkg/expression"
	myJson "MicroserviceCalculatorProject/orchestrator/pkg/json"
	"strings"

	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/streadway/amqp"
)

func (apiCfg *apiConfig) handlerProcessExpression(w http.ResponseWriter, r *http.Request) {
	type clientRequest struct {
		Expression string `json:"expression"`
	}

	decoder := json.NewDecoder(r.Body)

	request := clientRequest{}
	err := decoder.Decode(&request)

	if err != nil {
		myJson.RespondWithError(w, 400, err.Error())
		return
	}

	expr := expression.FormatExpression(request.Expression)

	if expression.IsValid(expr) {
		expressionFromDB, err := apiCfg.DB.GetExpressionByID(r.Context(), expression.CreateIdempotentKey(expr))

		if err != nil && err.Error() == "sql: no rows in result set" {
			status, subexpressions, err := registrateNewExpression(apiCfg, r, expr)

			if err != nil {
				myJson.RespondWithJSON(w, status, fmt.Sprintf("Error: %v", err))
				return
			}
			myJson.RespondWithJSON(w, status, "It's new expression")

			tasks := prepareSubexpressionsForSend(subexpressions)

			subexpressionsDurations, err := apiCfg.DB.GetDurations(r.Context())
			if err != nil {
				myJson.RespondWithError(w, 500, "Couldn't get subexpressionsDurations")
				return
			}
			sendSubexpressions(tasks, subexpressionsDurations)

		} else if err != nil {
			myJson.RespondWithError(w, 400, fmt.Sprintf("Couldn't get expression: %v", err))
		} else {
			myJson.RespondWithJSON(w, 200, expressionFromDB)
		}

	} else {
		myJson.RespondWithError(w, 400, "invalid expression")
	}
}

func (apiCfg *apiConfig) handlerGetExpression(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Expression ID: " + id))
}

//-------
//logical segments
//-------

func registrateNewExpression(apiCfg *apiConfig, r *http.Request, expr []string) (int, []database.Subexpression, error) {
	subexpressionMap := make(map[int]string)
	subexpressions := make([]database.Subexpression, 0)

	err := expression.ProcessExpression(expr, subexpressionMap)
	if err != nil {
		return 400, nil, err
	}

	newExpression, err := apiCfg.DB.CreateExpression(r.Context(), database.CreateExpressionParams{
		ID:                   expression.CreateIdempotentKey(expr),
		ExpressionBody:       strings.Join(expr, " "),
		CountOfSubexpression: int32(len(subexpressionMap)),
	})

	if err != nil {
		return 500, nil, err
	}

	for number, body := range subexpressionMap {
		subexpression, err := apiCfg.DB.CreateSubexpression(r.Context(), database.CreateSubexpressionParams{
			ExpressionID:        newExpression.ID,
			SubexpressionNumber: int32(number),
			SubexpressionBody:   body,
		})
		if err != nil {
			return 500, nil, err
		}
		subexpressions = append(subexpressions, subexpression)
	}

	return 200, subexpressions, nil
}

func sendSubexpressions(tasks []collection.AgentsTask, durations []database.OperatorsDuration) error {

	// convert database structs to simple map
	durationsMap := make(map[string]float64)
	for _, v := range durations {
		durationsMap[v.OperatorName] = v.OperatorDuration
	}

	//connect to RabbitMQ & send tasks
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
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
		return err
	}

	for _, task := range tasks {
		task.OperatorDuration = durationsMap[task.Operator]

		body, err := json.Marshal(task)
		if err != nil {
			return err
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
			return err
		}
	}

	return nil
}

func prepareSubexpressionsForSend(subexpressions []database.Subexpression) []collection.AgentsTask {
	agentsTasks := make([]collection.AgentsTask, 0)

	for _, subexpression := range subexpressions {
		specialOperands := expression.GetSubsexprNumbersBySubsexpr(subexpression.SubexpressionBody)

		for i := 0; i < len(specialOperands); i++ {
			for j := 0; j < len(subexpressions); j++ {
				if int(subexpressions[j].SubexpressionNumber) == specialOperands[i] && subexpressions[j].SubexpressionStatusID == 1 {
					subexpression.SubexpressionBody =
						strings.ReplaceAll(subexpression.SubexpressionBody,
							fmt.Sprintf("{%d}", specialOperands[i]),
							fmt.Sprint(subexpressions[j].SubexpressionResult))
				}
			}
		}

		if !expression.IsContainsUnknownVar(subexpression.SubexpressionBody) {
			agentsTasks = append(agentsTasks, expression.ConvertSubexpressionToAgentsTask(subexpression))
		}
	}

	return agentsTasks
}

func getSubexpressionsResults(queueName string) (pkg.Response, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return pkg.Response{}, err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return pkg.Response{}, err
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare(
		"TasksResults", // Queue name
		true,           // Durable
		false,          // Delete when unused
		false,          // Exclusive
		false,          // No-wait
		nil,            // Arguments
	)
	if err != nil {
		return pkg.Response{}, err
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
		return pkg.Response{}, err
	}

	d := <-msgs
	var message pkg.Response
	err = json.Unmarshal(d.Body, &message)

	if err != nil {
		return pkg.Response{}, err
	}

	err = d.Ack(false)

	return message, err
}
