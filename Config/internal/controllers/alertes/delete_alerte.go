package alertes

import (
	"encoding/json"
	"middleware/example/internal/helpers"
	"middleware/example/internal/services/alertes"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
)

// DeleteAlerte
// @Tags         alertes
// @Summary      Delete an alerte by ID.
// @Description  Delete an alerte by ID.
// @Param        id path string true "Alerte UUID formatted ID"
// @Success      200 {object} models.Alerte "Deleted alerte"
// @Failure      400 "Cannot parse id"
// @Failure      404 "Alerte not found"
// @Failure      500 "Something went wrong"
// @Router       /alertes/{id} [delete]

func DeleteAlerte(w http.ResponseWriter, r *http.Request) {
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

	deletedAlerte, err := alertes.DeleteAlerteById(id)
	if err != nil {
		body, status := helpers.RespondError(err)
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(deletedAlerte)
	_, _ = w.Write(body)
}
