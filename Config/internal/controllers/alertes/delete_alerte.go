package alertes

import (
	"middleware/example/internal/helpers"
	"middleware/example/internal/services/alertes"
	"net/http"

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
	ctx := r.Context()
	alerteId, _ := ctx.Value("alerteId").(uuid.UUID)

	err := alertes.DeleteAlerteById(alerteId)
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
