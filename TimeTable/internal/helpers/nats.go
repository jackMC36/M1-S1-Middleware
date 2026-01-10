package helpers

import (
	"time"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

var NatsConn *nats.Conn

func InitNats() error {
	// already initialized
	if NatsConn != nil && NatsConn.IsConnected() {
		return nil
	}

	nc, err := nats.Connect(
		nats.DefaultURL, // "nats://127.0.0.1:4222"
		nats.Name("Timetable service"),
		nats.Timeout(5*time.Second),
		nats.ReconnectWait(1*time.Second),
		nats.MaxReconnects(-1), // infinite reconnect
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			logrus.Warnf("NATS disconnected: %v", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			logrus.Infof("NATS reconnected to %s", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(_ *nats.Conn) {
			logrus.Warn("NATS connection closed")
		}),
	)
	if err != nil {
		return err
	}

	NatsConn = nc
	return nil
}

func CloseNats() {
	if NatsConn != nil {
		NatsConn.Drain()
		NatsConn.Close()
		NatsConn = nil
	}
}
