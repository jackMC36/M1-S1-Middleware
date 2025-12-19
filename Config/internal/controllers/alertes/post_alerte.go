package alertes

import (
	"encoding/json"
	"middleware/example/internal/helpers"
	"middleware/example/internal/services/alertes"
	"net/http"
)

func PostNewAlerte(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Email    string `json:"email"`
		AgendaId string `json:"agendaid"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	alerte, err := alertes.PostNewAlerte(requestData.AgendaId, requestData.Email)
	if err != nil {
		body, status := helpers.RespondError(err)
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	body, _ := json.Marshal(alerte)
	_, _ = w.Write(body)
}
