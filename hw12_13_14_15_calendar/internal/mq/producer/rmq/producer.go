package internalrmqproducer

import (
	"context"
	"log"
	"time"

	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/mq"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/mq/producer"
	"github.com/goccy/go-json"
	"github.com/streadway/amqp"
)

type Producer struct {
	Addr         string
	QueueName    string
	ExchangeName string
	connection   *amqp.Connection
	channel      *amqp.Channel
	logger       logger.Logger
}

func New(conf *producer.Config, logger logger.Logger) *Producer {
	addr, err := mq.GetAddress(conf.Protocol, conf.Host, conf.Port, conf.Username, conf.Password)
	if err != nil {
		log.Fatal(err)
	}
	return &Producer{
		Addr:         addr,
		ExchangeName: conf.ExchangeName,
		QueueName:    conf.QueueName,
		logger:       logger,
	}
}

func (p *Producer) Connect(ctx context.Context) error {
	var err error
	p.logger.Info("connect to rmq")
	p.connection, err = amqp.Dial(p.Addr)
	if err != nil {
		return err
	}

	p.channel, err = p.connection.Channel()
	if err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}

func (p *Producer) Close(ctx context.Context) error {
	err := p.channel.Close()
	if err != nil {
		return err
	}

	err = p.connection.Close()
	if err != nil {
		return err
	}
	p.logger.Info("rmq client shutdown successfully")
	return nil
}

func (p *Producer) Publish(ctx context.Context, data interface{}) error {
	encoded, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return p.channel.Publish(
		p.ExchangeName,
		p.QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        encoded,
			Timestamp:   time.Now(),
		},
	)
}
