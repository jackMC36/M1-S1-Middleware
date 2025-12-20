package helpers

import (
	"context"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

func EnsureStreams() error {
	js, err := jetstream.New(NatsConn)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Stream that receives scheduler messages and also the change events you publish.
	// Adjust subjects to your real ones.
	_, err = js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     "EventsStream",
		Subjects: []string{"Scheduler.Events", "Events.Changed"},
	})
	return err
}
