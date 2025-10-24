package agendas

import (
	"encoding/json"
	"middleware/example/internal/helpers"
	"middleware/example/internal/services/agendas"
	"net/http"

	"github.com/gofrs/uuid"
)

// GetAgenda
// @Tags         agendas
// @Summary      Get a Agenda.
// @Description  Get a Agenda.
// @Param        id           	path      string  true  "Agenda UUID formatted ID"
// @Success      200            {object}  models.Agenda
// @Failure      422            "Cannot parse id"
// @Failure      500            "Something went wrong"
// @Router       /agendas/{id} [get]
func GetAgenda(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	AgendaId, _ := ctx.Value("AgendaId").(uuid.UUID) // getting key set in context.go

	Agenda, err := agendas.GetAgendaById(AgendaId)
	if err != nil {
		body, status := helpers.RespondError(err)
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(Agenda)
	_, _ = w.Write(body)
	return
}
