package main

import (
	alertes_consumers "middleware/example/internal/consumers"
	"middleware/example/internal/helpers"

	"github.com/sirupsen/logrus"
)

func main() {
	err := helpers.InitNats()
	if err != nil {
		logrus.Fatalf("failed to connect to NATS: %v", err)
	}
	defer helpers.CloseNats()

	consumer, err := alertes_consumers.AlertConsumer()
	if err != nil {
		logrus.Warnf("error during nats consumer creation : %v", err)
	} else {
		err = alertes_consumers.Consume(*consumer)
		if err != nil {
			logrus.Warnf("error during nats consume : %v", err)
		}
	}
}
