package consumers

import (
	"context"
	"fmt"
	"time"

	"middleware/example/internal/helpers"
	"middleware/example/internal/services"

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

	stream, err := js.Stream(ctx, "ALERTES")
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

		err := services.SendMail(
			"jacques.KOZIK@etu.uca.fr", // TODO: get from message or config
			"Alert: "+msg.Subject(),
			string(msg.Data()),
		)
		if err != nil {
			logrus.Warnf("Failed to send email alert: %v", err)
		} else {
			logrus.Info("Email alert sent successfully")
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
