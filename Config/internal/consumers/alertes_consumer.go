package consumers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"middleware/example/internal/helpers"
	"middleware/example/internal/models"
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

	stream, err := js.Stream(ctx, "EVENT_CHANGES")
	if err != nil {
		return nil, fmt.Errorf("failed to get stream: %w", err)
	}

	consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:       "alerter-consumer",
		AckPolicy:     jetstream.AckExplicitPolicy,
		FilterSubject: "Events.Changed",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	return &consumer, nil
}

func Consume(consumer jetstream.Consumer) (err error) {
	cc, err := consumer.Consume(func(msg jetstream.Msg) {
		logrus.Infof("Received event on %s", msg.Subject())

var alert models.TimetableAlert
		if err := json.Unmarshal(msg.Data(), &alert); err != nil {
			logrus.Warnf("Failed to unmarshal payload into TimetableAlert: %v", err)
			_ = msg.Ack()
			return
		}

		agendaID := alert.AgendaID.String()

		alertes, svcErr := alertesService.GetAlertesByAgendaId(agendaID)
		if svcErr != nil {
			logrus.Warnf("Failed to fetch alertes for agenda %s: %v", agendaID, svcErr)
			_ = msg.Ack()
			return
		}

		if len(alert.Changes) == 0 {
			_ = msg.Ack()
			return
		}

	
		for _, a := range alertes {
			if a.Email == "" {
				continue
			}
			subject := "Event updated: " + alert.UID
			body := fmt.Sprintf("Event %s updated for agenda %s\n\nChanges:\n", alert.UID, alert.AgendaID.String())
			for _, c := range alert.Changes {
				body += fmt.Sprintf(" - %s: \"%s\" -> \"%s\"\n", c.Field, c.Before, c.After)
			}
			body += "\nNew event:\n" + prettyJSON(alert.After)

			if err := services.SendMail(a.Email, subject, body); err != nil {
				logrus.Warnf("Failed to send email to %s: %v", a.Email, err)
			} else {
				logrus.Infof("Email sent to %s", a.Email)
			}
		}

		_ = msg.Ack()
	})
	if err != nil {
		return fmt.Errorf("error creating consumer: %w", err)
	}

	<-cc.Closed()
	cc.Stop()

	return err
}

func prettyJSON(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(b)
}
