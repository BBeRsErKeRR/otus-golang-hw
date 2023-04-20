package rmq

import (
	"github.com/streadway/amqp"
)

type MessageQueue struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func (m *MessageQueue) Connect(addr string) error {
	var err error
	m.Connection, err = amqp.Dial(addr)
	if err != nil {
		return err
	}

	m.Channel, err = m.Connection.Channel()
	if err != nil {
		return err
	}
	return nil
}

func (m *MessageQueue) Close() error {
	err := m.Channel.Close()
	if err != nil {
		return err
	}
	err = m.Connection.Close()
	if err != nil {
		return err
	}
	return nil
}
