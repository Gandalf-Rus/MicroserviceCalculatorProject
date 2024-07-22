package orchestrator

import (
	"MicroserviceCalculatorProject/internal/database"
	myJson "MicroserviceCalculatorProject/pkg/json"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

func (apiConfig *apiConfig) handlerPostDuration(w http.ResponseWriter, r *http.Request) {
	type clientRequest struct {
		Operator string  `json:"operator"`
		Duration float64 `json:"duration"`
	}

	decoder := json.NewDecoder(r.Body)

	request := clientRequest{}
	err := decoder.Decode(&request)

	if err != nil {
		myJson.RespondWithError(w, 400, err.Error())
	}

	_, err = apiConfig.DB.EditDuration(r.Context(), database.EditDurationParams{
		OperatorName:     request.Operator,
		OperatorDuration: request.Duration,
	})
	if err != nil {
		myJson.RespondWithError(w, 400, err.Error())
	}

	myJson.RespondWithJSON(w, 200, struct{}{})

}

func (apiCfg *apiConfig) handlerGetDuration(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	duration, err := apiCfg.DB.GetDurationByName(r.Context(), name)

	if err != nil {
		myJson.RespondWithError(w, 400, err.Error())
	}
	myJson.RespondWithJSON(w, 200, duration)

}
