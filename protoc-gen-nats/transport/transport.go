package transport

import (
	"time"

	"github.com/nats-io/nats.go"
)

type Transport interface {
	Request(subj string, data []byte, timeout time.Duration) (*nats.Msg, error)
}

type NatsTransport struct {
	nc *nats.Conn
}

func (t *NatsTransport) Request(subj string, data []byte, timeout time.Duration) (*nats.Msg, error) {
	return t.nc.Request(subj, data, timeout)
}

func NewNatsTransport(nc *nats.Conn) Transport {
	return &NatsTransport{
		nc: nc,
	}
}
