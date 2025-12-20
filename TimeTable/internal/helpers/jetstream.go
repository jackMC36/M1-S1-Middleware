package helpers

import (
	"context"
	"time"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
)

func EnsureStreams() error {
	js, err := jetstream.New(NatsConn)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

    // Scheduler owns Scheduler.> in stream "EVENTS"
    if _, err := js.Stream(ctx, "EVENTS"); err != nil {
        return fmt.Errorf(`required stream "EVENTS" not found - start scheduler first : %w`, err)
    }

	    _, err = js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
        Name:     "EVENT_CHANGES",
        Subjects: []string{"Events.>"},
    })

    return nil
}
