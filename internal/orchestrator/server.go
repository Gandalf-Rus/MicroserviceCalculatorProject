package orchestrator

import (
	"MicroserviceCalculatorProject/internal/database"
	"MicroserviceCalculatorProject/internal/entities"
	"MicroserviceCalculatorProject/pkg/logger"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"

	_ "github.com/lib/pq"
)

type Server struct {
	*http.Server
	apiCfg apiConfig
}

type apiConfig struct {
	DB *database.Queries
}

func New() Server {

	//setup logger
	myLogger := logger.SetupInfoLogger()

	godotenv.Load("../.env")

	portString := os.Getenv("PORT")
	if portString == "" {
		myLogger.Fatal("PORT is not found in envoriment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		myLogger.Fatal("DB_URL is not found in envoriment")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		myLogger.Fatalf("Can't connect to database: %v", err)
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
	}

	router := chi.NewRouter()

	router.Use(logger.LoggingMiddleware(myLogger))

	apiRouter := chi.NewRouter()
	apiRouter.Get("/", handlerMainMenu)
	apiRouter.Get("/expression/{id}", apiCfg.handlerGetExpression)
	apiRouter.Post("/expression", apiCfg.handlerProcessExpression)
	apiRouter.Post("/durations", apiCfg.handlerPostDuration)
	apiRouter.Get("/durations/{name}", apiCfg.handlerGetDuration)

	router.Mount("/api", apiRouter)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port: %v", portString)

	return Server{srv, apiCfg}
}

func (server Server) Run() {
	//load all unresolve subexpression to RabbitMQ
	go func() {
		unresolveSubexpression, err := server.apiCfg.DB.GetSubexpressionByStatusID(context.Background(), 2)
		if err != nil {
			logger.SetupInfoLogger().Printf("database quere error: %v", err)
		}
		subexpressionsDurations, err := server.apiCfg.DB.GetDurations(context.Background())
		if err != nil {
			logger.SetupInfoLogger().Printf("database quere error: %v", err)
		}

		tasks := prepareSubexpressionsForSend(&server.apiCfg, unresolveSubexpression)

		sendSubexpressions(tasks, subexpressionsDurations)
	}()

	//check queue (RabbitMQ) on subexpressions results & send new unresolves subexpression
	go func() {
		conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		if err != nil {
			log.Printf("Error: %v", err)
		}
		defer conn.Close()

		ch, err := conn.Channel()
		if err != nil {
			log.Printf("Error: %v", err)
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
			log.Printf("Error: %v", err)
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
			log.Printf("Error: %v", err)
		}

		for d := range msgs {

			var subexprResult entities.AgentResponse
			err = json.Unmarshal(d.Body, &subexprResult)
			if err != nil {
				log.Printf("Error: %v", err)
			}
			err = d.Ack(false)
			if err != nil {
				log.Printf("Error: %v", err)
			}

			logger.SetupInfoLogger().Printf("server: %v", subexprResult)

			server.apiCfg.DB.EditSubexpressions(context.Background(), database.EditSubexpressionsParams{
				ExpressionID:          subexprResult.ExpressionID,
				SubexpressionNumber:   int32(subexprResult.SubexpressionNumber),
				SubexpressionStatusID: int32(subexprResult.StatusID),
				SubexpressionResult: sql.NullFloat64{
					Float64: subexprResult.Result,
					Valid:   true},
			})

			if !EditExpressionIfExpressionReady(&server.apiCfg, subexprResult.ExpressionID) {

				unresolveSubexpression, err := server.apiCfg.DB.GetSubexpressionByExprID(context.Background(), subexprResult.ExpressionID)
				if err != nil {
					logger.SetupInfoLogger().Printf("database quere error: %v", err)
				}

				subexpressionsDurations, err := server.apiCfg.DB.GetDurations(context.Background())
				if err != nil {
					logger.SetupInfoLogger().Printf("database quere error: %v", err)
				}
				tasks := prepareSubexpressionsForSend(&server.apiCfg, unresolveSubexpression)
				sendSubexpressions(tasks, subexpressionsDurations)
			}
		}
	}()

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
