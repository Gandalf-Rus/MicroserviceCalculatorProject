package server

import "net/http"

func handlerMainMenu(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hi there! ```Walcome to Distributed Calculator API```\nFor start enter commad:\n"))
	w.Write([]byte("\t/expression/{id} GET request - get expression result/status\n"))
	w.Write([]byte("\t/durations GET request - get operator duration\n"))
	w.Write([]byte("\t/durations/{name} POST request - send operator duration\n"))
}
