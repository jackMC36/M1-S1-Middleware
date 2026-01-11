package alertes

import (
	"encoding/json"
	"middleware/example/internal/helpers"
	"middleware/example/internal/services/alertes"
	"net/http"

	"github.com/gofrs/uuid"
)

// Getalerte
// @Tags         alertes
// @Summary      Get a alerte.
// @Description  Get a alerte.
// @Param        id           	path      string  true  "alerte UUID formatted ID"
// @Success      200            {object}  models.Alerte
// @Failure      422            "Cannot parse id"
// @Failure      500            "Something went wrong"
// @Router       /alertes/{id} [get]
func GetAlerte(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	alerteId, _ := ctx.Value("alerteId").(uuid.UUID)

	alerte, err := alertes.GetalerteById(alerteId)
	if err != nil {
		body, status := helpers.RespondError(err)
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(alerte)
	_, _ = w.Write(body)
}
