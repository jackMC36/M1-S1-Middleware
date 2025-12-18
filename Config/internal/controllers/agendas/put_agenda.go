package agendas

import (
	"encoding/json"
	"middleware/example/internal/helpers"
	"middleware/example/internal/services/agendas"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
)

// PutAgenda
// @Tags         agendas
// @Summary      Update an agenda by ID.
// @Description  Update an agenda by ID.
// @Accept       json
// @Produce      json
// @Param        id path string true "Agenda UUID formatted ID"
// @Param        agenda body models.Agenda true "Updated agenda object"
// @Success      200 {object} models.Agenda
// @Failure      400 "Bad Request"
// @Failure      404 "Agenda not found"
// @Failure      500 "Something went wrong"
// @Router       /agendas/{id} [put]
func PutAgenda(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		body, _ := json.Marshal(map[string]string{"error": "ID parameter is required"})
		_, _ = w.Write(body)
		return
	}

	id, err := uuid.FromString(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		body, _ := json.Marshal(map[string]string{"error": "Invalid ID format"})
		_, _ = w.Write(body)
		return
	}

	var requestData struct {
		Name  string `json:"name"`
		UCAID string `json:"ucaid"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		body, _ := json.Marshal(map[string]string{"error": "Invalid JSON"})
		_, _ = w.Write(body)
		return
	}

	updatedAgenda, err := agendas.UpdateAgendaById(id, requestData.Name, requestData.UCAID)
	if err != nil {
		body, status := helpers.RespondError(err)
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(updatedAgenda)
	_, _ = w.Write(body)
}
