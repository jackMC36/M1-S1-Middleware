package alertes

import (
	"encoding/json"
	"middleware/example/internal/helpers"
	"middleware/example/internal/services/alertes"
	"net/http"
)

// Getalertes
// @Tags         alertes
// @Summary      Get all alertes.
// @Description  Get all alertes.
// @Success      200            {array}  models.User
// @Failure      500             "Something went wrong"
// @Router       /alertes [get]
func GetAlertes(w http.ResponseWriter, _ *http.Request) {
	alertes, err := alertes.GetAllAlertes()
	if err != nil {
		body, status := helpers.RespondError(err)
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(alertes)
	_, _ = w.Write(body)
	return
}
