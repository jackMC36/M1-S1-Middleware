package agendas

import (
	"encoding/json"
	"middleware/example/internal/helpers"
	"middleware/example/internal/services/agendas"
	"net/http"

	"github.com/gofrs/uuid"
)

func PostNewAgenda(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		ID    uuid.UUID `json:"id"`
		Name  string    `json:"name"`
		UCAID uuid.UUID `json:"ucaid"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	agenda, err := agendas.PostNewAgenda(requestData.ID, requestData.Name, requestData.UCAID)
	if err != nil {
		body, status := helpers.RespondError(err)
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	body, _ := json.Marshal(agenda)
	_, _ = w.Write(body)
}
