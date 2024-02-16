package server

import (
	"MicroserviceCalculatorProject/orchestrator/internal/database"
	"MicroserviceCalculatorProject/orchestrator/pkg/expression"
	myJson "MicroserviceCalculatorProject/orchestrator/pkg/json"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
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

	exp := expression.FormatExpression(request.Expression)

	if expression.IsValid(exp) {
		subexpressionMap := make(map[int]string)
		err = expression.ProcessExpression(exp, subexpressionMap)
		if err != nil {
			myJson.RespondWithError(w, 400, err.Error())
			return
		}

		expression, err := apiCfg.DB.CreateExpression(r.Context(), database.CreateExpressionParams{
			ID:                 expression.CreateIdempotentKey(exp),
			ExpressionBody:     strings.Join(exp, " "),
			ExpressionStatusID: 2,
			ExpressionResult: sql.NullFloat64{
				Valid: false,
			},
		})

		if err != nil {
			myJson.RespondWithError(w, 400, fmt.Sprintf("Couldn't create expression: %v", err))
			return
		}

		myJson.RespondWithJSON(w, 200, expression)

	} else {
		myJson.RespondWithError(w, 400, "invalid expression")
	}
}

func (apiCfg *apiConfig) handlerGetExpression(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Expression ID: " + id))
}
