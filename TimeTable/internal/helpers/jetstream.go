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
        return fmt.Errorf("cannot initialize JetStream: %w", err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if _, err := js.Stream(ctx, "EVENTS"); err != nil {
        return fmt.Errorf(
            `required stream "EVENTS" not found. Start the Scheduler to create the stream and ensure NATS is running with JetStream enabled: %w`,
            err,
        )
    }

    _, err = js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
        Name:     "EVENT_CHANGES",
        Subjects: []string{"Events.>"},
    })
    return err
}

