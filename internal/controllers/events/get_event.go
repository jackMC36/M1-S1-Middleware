package events

import (
	"encoding/json"
	"middleware/example/internal/helpers"
	"middleware/example/internal/services/events"
	"net/http"

	"github.com/gofrs/uuid"
)

// Getevent
// @Tags         events
// @Summary      Get a event.
// @Description  Get a event.
// @Param        id           	path      string  true  "event UUID formatted ID"
// @Success      200            {object}  models.event
// @Failure      422            "Cannot parse id"
// @Failure      500            "Something went wrong"
// @Router       /events/{id} [get]
func Getevent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	eventId, _ := ctx.Value("eventId").(uuid.UUID) // getting key set in context.go

	event, err := events.GetEventById(eventId)
	if err != nil {
		body, status := helpers.RespondError(err)
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(event)
	_, _ = w.Write(body)
	return
}
