package main

import (
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	internalrmqconsumer "github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/mq/consumer/rmq"
)

type Config struct {
	Logger   *logger.Config              `mapstructure:"logger"`
	Consumer *internalrmqconsumer.Config `mapstructure:"rmq"`
}

func NewConfig(configFile string) (Config, error) {
	conf := Config{}
	_, err := config.ReadConfigFile(configFile, "toml", &conf)
	return conf, err
}
