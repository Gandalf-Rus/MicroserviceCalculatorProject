package orchestrator

import "net/http"

func handlerMainMenu(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hi there! Walcome to ```Distributed Calculator API```\nFor start enter commad:\n\n"))

	w.Write([]byte("\tcurl localhost:8080/api/expression/{id} - get expression result/status\n\n"))

	w.Write([]byte("\tcurl localhost:8080/api/durations/{operatorName} - get operator duration\n(for example: curl localhost:8080/api/durations/*)\n\n"))

	w.Write([]byte("Post request look into Readme in `Примеры и запросы`\n\n"))
}
