package helpers

import (
	"log"

	"github.com/nats-io/nats.go"
)

var NatsConn *nats.Conn

func InitNats() error {
	var err error
	NatsConn, err = nats.Connect(nats.DefaultURL)
	if err != nil {
		return err
	}
	log.Println("Connected to NATS")
	return nil
}

func CloseNats() {
	if NatsConn != nil {
		NatsConn.Close()
	}
}
