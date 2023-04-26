package rmq

import (
	"errors"
	"fmt"

	pkgnet "github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/pkg/net"
	"github.com/streadway/amqp"
)

var ErrorUnsupportedProtocol = errors.New("invalid protocol")

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

func GetMqAddress(protocolArg, hostArg, portArg, userArg, passArg string) (string, error) {
	var address string

	if protocolArg != "amqp" {
		return address, ErrorUnsupportedProtocol
	}

	netAddress, err := pkgnet.GetAddress(hostArg, portArg)
	if err != nil {
		return address, err
	}

	if userArg != "" && passArg != "" {
		address = fmt.Sprintf("%v://%v:%v@%v", protocolArg, userArg, passArg, netAddress)
	} else {
		address = fmt.Sprintf("%v://%v", protocolArg, netAddress)
	}

	return address, nil
}
