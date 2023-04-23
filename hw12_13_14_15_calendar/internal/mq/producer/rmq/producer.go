package internalrmqproducer

import (
	"context"
	"log"
	"time"

	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/utils"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/pkg/rmq"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type Config struct {
	Host         string `mapstructure:"host"`
	Port         string `mapstructure:"port"`
	Protocol     string `mapstructure:"protocol"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	QueueName    string `mapstructure:"queue"`
	ExchangeName string `mapstructure:"exchange"`
}

type Producer struct {
	Addr         string
	QueueName    string
	ExchangeName string
	mq           rmq.MessageQueue
	logger       logger.Logger
}

func New(conf *Config, logger logger.Logger) *Producer {
	addr, err := utils.GetMqAddress(conf.Protocol, conf.Host, conf.Port, conf.Username, conf.Password)
	if err != nil {
		log.Fatal(err)
	}
	return &Producer{
		Addr:         addr,
		ExchangeName: conf.ExchangeName,
		QueueName:    conf.QueueName,
		logger:       logger,
		mq:           rmq.MessageQueue{},
	}
}

func (p *Producer) Connect(ctx context.Context) error {
	p.logger.Info("connect to rmq")
	err := p.mq.Connect(p.Addr)
	if err != nil {
		return err
	}
	return p.Declare()
}

func (p *Producer) Declare() error {
	err := p.mq.Channel.ExchangeDeclare(
		p.ExchangeName, // name
		"direct",       // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		p.logger.Error("failed to declare a queue", zap.Error(err))
	}

	_, err = p.mq.Channel.QueueDeclare(
		p.QueueName, // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		p.logger.Error("failed to declare a queue", zap.Error(err))
	}

	err = p.mq.Channel.QueueBind(
		p.QueueName,
		p.QueueName,
		p.ExchangeName,
		false,
		nil)

	return err
}

func (p *Producer) Close(ctx context.Context) error {
	err := p.mq.Close()
	if err != nil {
		return err
	}
	p.logger.Info("rmq client shutdown successfully")
	return nil
}

func (p *Producer) Publish(ctx context.Context, data []byte) error {
	return p.mq.Channel.Publish(
		p.ExchangeName,
		p.QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
			Timestamp:   time.Now(),
		},
	)
}
