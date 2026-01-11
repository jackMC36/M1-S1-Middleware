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
	alerttpl "middleware/example/internal/templates/alertes"


	"github.com/nats-io/nats.go/jetstream"
	"github.com/sirupsen/logrus"
)

func AlertConsumer() (jetstream.Consumer, error) {
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

	return consumer, nil
}

func Consume(consumer jetstream.Consumer) (err error) {
	cc, err := consumer.Consume(func(msg jetstream.Msg) {
		logrus.Infof("Received event on %s", msg.Subject())

		var alert models.TimetableAlert
		if err := json.Unmarshal(msg.Data(), &alert); err != nil {
			logrus.Warnf("Failed to unmarshal payload into TimetableAlert: %v", err)
			_ = msg.Term()
			return
		}

		if len(alert.Changes) == 0 {
			_ = msg.Ack()
			return
		}
		
		agendaID := alert.AgendaID.String()

		alertes, svcErr := alertesService.GetAlertesByAgendaId(agendaID)
		if svcErr != nil {
			logrus.Warnf("Config API error for agenda %s: %v", agendaID, svcErr)
			_ = msg.Nak()
			return
		}
			
		loc, _ := time.LoadLocation("Europe/Paris")

		tplData := struct {
			EventName   string
			Start       string
			End         string
			Location    string
			Description string
			Changes     []helpers.ChangeLine
			DebugRef    string
		}{
			EventName:   alert.After.Name,
			Start:       alert.After.Start.In(loc).Format("02/01/2006 15:04"),
			End:         alert.After.End.In(loc).Format("02/01/2006 15:04"),
			Location:    alert.After.Location,
			Description: alert.After.Description,
			Changes:     helpers.MapChanges(alert.Changes),
		}

		body, matter, tplErr := alerttpl.GetStringFromEmbeddedTemplate("event_changed.txt.tmpl", tplData)

		if tplErr != nil {
			logrus.Warnf("Template render error: %v", tplErr)
			_ = msg.Term()
			return
		}

		subject := matter.Subject
		if subject == "" {
			subject = "EDU Middleware - modification d'un cours"
		}

		for _, a := range alertes {
			if a.Email == "" {
				continue
			}

			if err := services.SendMail(a.Email, subject, body); err != nil {
				logrus.Warnf("Failed to send email to %s: %v", a.Email, err)
				_ = msg.Nak()
				return
			} else {
				logrus.Infof("Email sent to %s", a.Email)
			}
		}

		_ = msg.Ack()
	})
	if err != nil {
		return fmt.Errorf("error creating consumer: %w", err)
	}

	defer cc.Stop()
	<-cc.Closed()

	return err
}
