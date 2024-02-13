package server

import (
	"MicroserviceCalculatorProject/orchestrator/pkg/expression"
	"MicroserviceCalculatorProject/orchestrator/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Orchestrator struct {
	mux *mux.Router
}

func New() Orchestrator {
	serveMux := mux.NewRouter()

	//setup logger
	myLogger := logger.SetupLogger()
	serveMux.Use(logger.LoggingMiddleware(myLogger))

	//handler for main menu
	serveMux.HandleFunc("/api", mainMenu)

	// Обработчик POST запросов на /expression
	serveMux.HandleFunc("/api/expression", postExpression).Methods("POST")

	// Обработчик GET запросов на /expression/{id}
	serveMux.HandleFunc("/api/expression/{id}", getExpressionByID).Methods("GET")

	// Use router Gorilla Mux
	http.Handle("/", serveMux)

	return Orchestrator{serveMux}
}

func (orchestrator Orchestrator) Run() (func(context.Context) error, error) {

	srv := &http.Server{Addr: ":8080", Handler: orchestrator.mux}

	// start server
	if err := srv.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
	// вернем функцию для завершения работы сервера
	return srv.Shutdown, nil
}

//-------
//functions for handle paths
//-------

type Request struct {
	Expression string `json:"expression"`
}

func mainMenu(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hi there! ''You use Distributed Calculator API''\nFor start enter commad:\n"))
}

func postExpression(w http.ResponseWriter, r *http.Request) {

	request := Request{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}

	exp := expression.FormatExpression(request.Expression)

	fmt.Fprintf(w, "%s\n", exp)

	if expression.IsValid(exp) {
		subexpressionMap := make(map[int]string)
		err = expression.ProcessExpression(exp, subexpressionMap)
		if err != nil {
			w.WriteHeader(400)
		}

		for k, v := range subexpressionMap {
			w.Write([]byte(fmt.Sprintf("%d:\t", k) + v + "\n"))
		}

		w.Write([]byte(expression.CreateIdempotentKey(exp)))

	} else {
		w.Write([]byte("invalid expression"))
	}
}

func getExpressionByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Fprintf(w, "ID: %s\n", id)
}
