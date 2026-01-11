package agendas

import (
	"encoding/json"
	"middleware/example/internal/helpers"
	"middleware/example/internal/services/agendas"
	"net/http"
)

// GetAgendas
// @Tags         agendas
// @Summary      Get all agendas
// @Description  Returns the list of agendas
// @Success      200  {array}   models.Agenda
// @Failure      500  {object}  helpers.ErrorResponse "Something went wrong"
// @Router       /agendas [get]

func GetAgendas(w http.ResponseWriter, _ *http.Request) {
	agendas, err := agendas.GetAllagendas()
	if err != nil {
		body, status := helpers.RespondError(err)
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(agendas)
	_, _ = w.Write(body)
}
