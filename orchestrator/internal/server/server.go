package server

import (
	"MicroserviceCalculatorProject/orchestrator/internal/database"
	"MicroserviceCalculatorProject/orchestrator/pkg/logger"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type Server struct {
	*http.Server
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

	router.Mount("/api", apiRouter)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port: %v", portString)

	return Server{srv}
}

func (server Server) Run() {

	// clear log file every 10 minute
	go func() {
		for {
			logger.ClearFileLog()
			<-time.After(time.Minute * 20)
		}
	}()

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
