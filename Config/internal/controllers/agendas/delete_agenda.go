package agendas

import (
	"encoding/json"
	"middleware/example/internal/helpers"
	"middleware/example/internal/services/agendas"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
)

// DeleteAgenda
// @Tags         agendas
// @Summary      Delete an agenda by ID.
// @Description  Delete an agenda by ID.
// @Param        id path string true "Agenda UUID formatted ID"
// @Success      200 {object} models.Agenda "Deleted agenda"
// @Failure      400 "Cannot parse id"
// @Failure      404 "Agenda not found"
// @Failure      500 "Something went wrong"
// @Router       /agendas/{id} [delete]
func DeleteAgenda(w http.ResponseWriter, r *http.Request) {
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

	err = agendas.DeleteAgendaById(id)
	if err != nil {
		body, status := helpers.RespondError(err)
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
