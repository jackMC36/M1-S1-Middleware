package alertes

import (
	"encoding/json"
	"middleware/example/internal/helpers"
	"middleware/example/internal/services/alertes"
	"net/http"

	"github.com/gofrs/uuid"
)

func PostNewAlerte(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Id       *uuid.UUID `json:"id"`
		Email    string     `json:"email"`
		AgendaId *uuid.UUID `json:"agendaid"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	alerte, err := alertes.PostNewAlerte(*requestData.AgendaId, requestData.Email)
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
	return
}
