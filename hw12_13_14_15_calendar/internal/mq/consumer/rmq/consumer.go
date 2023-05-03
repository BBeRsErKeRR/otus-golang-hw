package internalrmqconsumer

import (
	"context"
	"log"
	"time"

	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/pkg/rmq"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type Config struct {
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	Protocol        string `mapstructure:"protocol"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	Subscription    string `mapstructure:"subscription"`
	ConsumerName    string `mapstructure:"consumer_name"`
	ExchangeName    string `mapstructure:"exchange"`
	StatusQueueName string `mapstructure:"status_queue"`
}

type Consumer struct {
	Addr            string
	Subscription    string
	ConsumerName    string
	ExchangeName    string
	StatusQueueName string
	mq              rmq.MessageQueue
	logger          logger.Logger
	StatusChan      *amqp.Channel
}

func New(conf *Config, logger logger.Logger) *Consumer {
	addr, err := rmq.GetMqAddress(conf.Protocol, conf.Host, conf.Port, conf.Username, conf.Password)
	if err != nil {
		log.Fatal(err)
	}
	return &Consumer{
		Addr:            addr,
		Subscription:    conf.Subscription,
		ConsumerName:    conf.ConsumerName,
		ExchangeName:    conf.ExchangeName,
		StatusQueueName: conf.StatusQueueName,
		logger:          logger,
		mq:              rmq.MessageQueue{},
	}
}

func (c *Consumer) Connect(ctx context.Context) error {
	c.logger.Info("connect to rmq")
	err := c.mq.Connect(c.Addr)
	if err != nil {
		c.logger.Error("failed to connect", zap.Error(err))
		return err
	}

	c.StatusChan, err = c.mq.Connection.Channel()
	if err != nil {
		c.logger.Error("failed open status chan", zap.Error(err))
		return err
	}

	err = c.StatusChan.ExchangeDeclare(
		c.ExchangeName, // name
		"direct",       // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		c.logger.Error("failed to declare a queue", zap.Error(err))
		return err
	}

	_, err = c.StatusChan.QueueDeclare(
		c.StatusQueueName, // name
		false,             // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		c.logger.Error("failed to declare a queue", zap.Error(err))
		return err
	}

	err = c.StatusChan.QueueBind(
		c.StatusQueueName,
		c.StatusQueueName,
		c.ExchangeName,
		false,
		nil)

	return err
}

func (c *Consumer) Close(ctx context.Context) error {
	err := c.StatusChan.Close()
	if err != nil {
		return err
	}
	err = c.mq.Close()
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

func (c *Consumer) PublishStatus(ctx context.Context, data []byte) error {
	return c.StatusChan.Publish(
		c.ExchangeName,
		c.StatusQueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
			Timestamp:   time.Now(),
		},
	)
}
