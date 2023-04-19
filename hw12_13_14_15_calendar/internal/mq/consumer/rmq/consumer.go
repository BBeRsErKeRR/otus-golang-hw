package internalrmqconsumer

import (
	"context"
	"log"

	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/mq"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/mq/consumer"
	"github.com/streadway/amqp"
)

type Consumer struct {
	Addr         string
	Subscription string
	ConsumerName string
	connection   *amqp.Connection
	channel      *amqp.Channel
	logger       logger.Logger
}

func New(conf *consumer.Config, logger logger.Logger) *Consumer {
	addr, err := mq.GetAddress(conf.Protocol, conf.Host, conf.Port, conf.Username, conf.Password)
	if err != nil {
		log.Fatal(err)
	}
	return &Consumer{
		Addr:         addr,
		Subscription: conf.Subscription,
		ConsumerName: conf.ConsumerName,
		logger:       logger,
	}
}

func (c *Consumer) Connect(ctx context.Context) error {
	var err error
	c.logger.Info("connect to rmq")
	c.connection, err = amqp.Dial(c.Addr)
	if err != nil {
		return err
	}

	c.channel, err = c.connection.Channel()
	if err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}

func (c *Consumer) Close(ctx context.Context) error {
	err := c.channel.Close()
	if err != nil {
		return err
	}
	err = c.connection.Close()
	if err != nil {
		return err
	}
	c.logger.Info("rmq client shutdown successfully")
	return nil
}

func (c *Consumer) Consume(ctx context.Context, f func(ctx context.Context, msg []byte)) error {
	msgs, err := c.channel.Consume(
		c.Subscription,
		c.ConsumerName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-msgs:
			if !ok {
				return nil
			}
			c.logger.Info("receive msg from a queue")
			f(ctx, msg.Body)
		}
	}
}
