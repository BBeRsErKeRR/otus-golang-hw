package internalrmqconsumer

import (
	"context"
	"log"

	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/mq/consumer"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/utils"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/pkg/rmq"
)

type Consumer struct {
	Addr         string
	Subscription string
	ConsumerName string
	mq           rmq.MessageQueue
	logger       logger.Logger
}

func New(conf *consumer.Config, logger logger.Logger) *Consumer {
	addr, err := utils.GetMqAddress(conf.Protocol, conf.Host, conf.Port, conf.Username, conf.Password)
	if err != nil {
		log.Fatal(err)
	}
	return &Consumer{
		Addr:         addr,
		Subscription: conf.Subscription,
		ConsumerName: conf.ConsumerName,
		logger:       logger,
		mq:           rmq.MessageQueue{},
	}
}

func (c *Consumer) Connect(ctx context.Context) error {
	c.logger.Info("connect to rmq")
	return c.mq.Connect(c.Addr)
}

func (c *Consumer) Close(ctx context.Context) error {
	err := c.mq.Close()
	if err != nil {
		return err
	}
	c.logger.Info("rmq client shutdown successfully")
	return nil
}

func (c *Consumer) Consume(ctx context.Context, f func(ctx context.Context, msg []byte)) error {
	msgs, err := c.mq.Channel.Consume(
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
			return ctx.Err()
		case msg, ok := <-msgs:
			if !ok {
				return nil
			}
			c.logger.Info("receive msg from a queue")
			f(ctx, msg.Body)
		}
	}
}
