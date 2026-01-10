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
    js, err := jetstream.New(helpers.NatsConn)
    if err != nil {
        return nil, err
    }    
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
            // New event: store it, but DO NOT notify
            ev := payload.Event
            if ev.ID == uuid.Nil {
                ev.ID = uuid.Must(uuid.NewV4())
            }
            ev.AgendaIDs = agendaIDs

            if err := repository.InsertEvent(ev); err != nil {
                return err
            }
            return nil
        }
        return err
    }

    before := *existingEvent
    after := payload.Event
    after.ID = existingEvent.ID
    after.AgendaIDs = agendaIDs

    changes := eventDiff(before, after)
    if len(changes) == 0 {
        return nil
    }

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
    alert := models.TimetableAlert{
        AgendaID: payload.AgendaID,
        UID:      after.UID,
        Changes:  changes,
        Before:   &before,
        After:    after,
    }

    return publishTimetableAlert(js, "Events.Changed", alert)
}


func publishTimetableAlert(js jetstream.JetStream, subject string, alert models.TimetableAlert) error {
	data, err := json.Marshal(alert)
	if err != nil {
		return err
	}
	_, err = js.Publish(context.Background(), subject, data)
	return err
}

func eventDiff(before models.Event, after models.Event) []models.FieldChange {
    changed := []models.FieldChange{}

    if before.Name != after.Name {
        changed = append(changed, models.FieldChange{
            Field: "name",
            Before: before.Name,
            After: after.Name,
        })
    }
    if before.Description != after.Description {
        changed = append(changed, models.FieldChange{
            Field: "description",
            Before: before.Description,
            After: after.Description,
        })
    }
    if before.Location != after.Location {
        changed = append(changed, models.FieldChange{
            Field: "location",
            Before: before.Location,
            After: after.Location,
        })      
    }
    if !before.Start.UTC().Equal(after.Start.UTC()) {
        changed = append(changed, models.FieldChange{
            Field:  "start",
            Before: before.Start.UTC().Format(time.RFC3339),
            After:  after.Start.UTC().Format(time.RFC3339),
        })
    }
    if !before.End.UTC().Equal(after.End.UTC()) {
        changed = append(changed, models.FieldChange{
            Field:  "end",
            Before: before.End.UTC().Format(time.RFC3339),
            After:  after.End.UTC().Format(time.RFC3339),
        })
    }

    return changed
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
            _ = msg.Nak()
            return
        }

        _ = msg.Ack()
    })

    <-cc.Closed()
    cc.Stop()
    return err
}




