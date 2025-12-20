package events

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/sirupsen/logrus"
	"middleware/example/internal/helpers"
	"middleware/example/internal/models"
	repository "middleware/example/internal/repositories/events"
)

func EventConsumer() (*jetstream.Consumer, error) {
    js, _ := jetstream.New(helpers.NatsConn)
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    stream, err := js.Stream(ctx, "EVENTS")
    if err != nil {
        return nil, err
    }


    consumer, err := stream.Consumer(ctx, "TimeTable_Consumer")
    if err != nil {
        consumer, err = stream.CreateConsumer(ctx, jetstream.ConsumerConfig{
            Durable:     "TimeTable_Consumer",
            Name:        "TimeTable_Consumer",
            Description: "Consumes timetable update events",
            FilterSubject: "Scheduler.Events",    
        })
        if err != nil {
            return nil, err
        }
        logrus.Infof("Created consumer")
    } else {
        logrus.Infof("Got existing consumer")
    }

    return &consumer, nil
}

func HandlePullMessage(js jetstream.JetStream, payload models.SchedulerPayload) error {
    agendaIDs := []uuid.UUID{payload.AgendaID}

    existingEvent, err := repository.GetEventByUID(payload.Event.UID)
    if err != nil {
        if err == sql.ErrNoRows {
            // New event
            ev := payload.Event
            if ev.ID == uuid.Nil {
                ev.ID = uuid.Must(uuid.NewV4())
            }
            ev.AgendaIDs = agendaIDs

            if err := repository.InsertEvent(ev); err != nil {
                return err
            }

            change := models.EventChange{
                Before: nil,
                After:  ev,
            }
            return publishEventChange(js, "Events.Changed", change)
        }
        return err
    }

    before := *existingEvent

    _, err = repository.UpdateEventByUID(
        payload.Event.UID,
        agendaIDs,
        payload.Event.Description,
        payload.Event.Name,
        payload.Event.Start,
        payload.Event.End,
        payload.Event.Location,
        payload.Event.LastUpdate,
    )
    if err != nil {
        return err
    }

    if err := repository.EventAgendasLink(existingEvent.ID, agendaIDs); err != nil {
        return err
    }

    // After
    after := payload.Event
    after.ID = existingEvent.ID
    after.AgendaIDs = agendaIDs

    change := models.EventChange{
        Before: &before,
        After:  after,
    }
    return publishEventChange(js, "Events.Changed", change)
}


func publishEventChange(js jetstream.JetStream, subject string, change models.EventChange) error {
	data, err := json.Marshal(change)
	if err != nil {
		return err
	}
	_, err = js.Publish(context.Background(), subject, data)
	return err
}




func Consume(consumer jetstream.Consumer) error {
    js, err := jetstream.New(helpers.NatsConn)
    if err != nil {
        return err
    }

    cc, err := consumer.Consume(func(msg jetstream.Msg) {
        var payload models.SchedulerPayload
        if err := json.Unmarshal(msg.Data(), &payload); err != nil {
            logrus.Errorf("invalid message (json): %v", err)
            _ = msg.Ack()
            return
        }

        if err := HandlePullMessage(js, payload); err != nil {
            logrus.Errorf("handle message failed: %v", err)
            _ = msg.Ack()
            return
        }

        _ = msg.Ack()
    })

    <-cc.Closed()
    cc.Stop()
    return err
}




