package consumers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"middleware/example/internal/helpers"
	"middleware/example/internal/services"
	alertesService "middleware/example/internal/services/alertes"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/sirupsen/logrus"
)

func AlertConsumer() (*jetstream.Consumer, error) {
	js, err := jetstream.New(helpers.NatsConn)
	if err != nil {
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stream, err := js.Stream(ctx, "EVENTS")
	if err != nil {
		return nil, fmt.Errorf("failed to get stream: %w", err)
	}

	consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:   "alerter-consumer",
		AckPolicy: jetstream.AckExplicitPolicy,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	return &consumer, nil
}

func Consume(consumer jetstream.Consumer) (err error) {
	cc, err := consumer.Consume(func(msg jetstream.Msg) {
		logrus.Infof("Received event on %s: %s", msg.Subject(), string(msg.Data()))

		// Parse payload to get agenda id
		var payload struct {
			AgendaID string `json:"agendaId"`
		}
		if err := json.Unmarshal(msg.Data(), &payload); err != nil {
			logrus.Warnf("Failed to unmarshal payload: %v", err)
		} else {
			// fetch alertes for this agenda
			alertes, svcErr := alertesService.GetAlertesByAgendaId(payload.AgendaID)
			if svcErr != nil {
				logrus.Warnf("Failed to fetch alertes for agenda %s: %v", payload.AgendaID, svcErr)
			} else {
				for _, a := range alertes {
					if a.Email == "" {
						continue
					}
					if err := services.SendMail(a.Email, "Alert: "+msg.Subject(), string(msg.Data())); err != nil {
						logrus.Warnf("Failed to send email to %s: %v", a.Email, err)
					} else {
						logrus.Infof("Email sent to %s", a.Email)
					}
				}
			}
		}

		if err := msg.Ack(); err != nil {
			logrus.Warnf("Error acknowledging message: %v", err)
		}
	})
	if err != nil {
		return fmt.Errorf("error creating consumer: %w", err)
	}

	<-cc.Closed()
	cc.Stop()

	return err
}
