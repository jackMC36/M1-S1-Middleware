package alertes

import (
	"encoding/json"
	"middleware/example/internal/helpers"
	"middleware/example/internal/services/alertes"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
)

// PutAlerte
// @Tags         alertes
// @Summary      Update an alerte by ID.
// @Description  Update an alerte by ID.
// @Accept       json
// @Produce      json
// @Param        id path string true "Alerte UUID formatted ID"
// @Param        alerte body models.Alerte true "Updated alerte object"
// @Success      200 {object} models.Alerte
// @Failure      400 "Bad Request"
// @Failure      404 "Alerte not found"
// @Failure      500 "Something went wrong"
// @Router       /alertes/{id} [put]
func PutAlerte(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL parameter
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		body, _ := json.Marshal(map[string]string{"error": "ID parameter is required"})
		_, _ = w.Write(body)
		return
	}

	// Parse UUID
	id, err := uuid.FromString(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		body, _ := json.Marshal(map[string]string{"error": "Invalid ID format"})
		_, _ = w.Write(body)
		return
	}

	var requestData struct {
		Id       uuid.UUID `json:"id"`
		Email    string    `json:"email"`
		AgendaId string    `json:"agendaid"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		body, _ := json.Marshal(map[string]string{"error": "Invalid JSON"})
		_, _ = w.Write(body)
		return
	}

	// Call service to update alerte
	updatedAlerte, err := alertes.UpdateAlerteById(id, requestData.Email, requestData.AgendaId)
	if err != nil {
		body, status := helpers.RespondError(err)
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(updatedAlerte)
	_, _ = w.Write(body)
}
