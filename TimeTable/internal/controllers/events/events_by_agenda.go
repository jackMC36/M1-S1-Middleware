package events

import (
	"encoding/json"
	"middleware/example/internal/helpers"
	"middleware/example/internal/services/events"
	"net/http"

	"github.com/gofrs/uuid"
)

// GetEventsByID
// @Tags         events
// @Summary      Get events for an agenda.
// @Description  Returns all events linked to the provided agenda ID.
// @Param        id            path      string  true  "Agenda UUID formatted ID"
// @Success      200            {array}   models.Event
// @Failure      422            "Cannot parse id"
// @Failure      500            "Something went wrong"
// @Router       /agendas/{id}/events [get]
func GetEventsByAgendaID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	agendaId, _ := ctx.Value("agendaId").(uuid.UUID) // getting key set in agendas.Context

	evts, err := events.GetEventsByAgenda(agendaId)
	if err != nil {
		body, status := helpers.RespondError(err)
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(evts)
	_, _ = w.Write(body)
	return
}
