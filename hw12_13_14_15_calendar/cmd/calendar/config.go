package main

import (
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/storage"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger *logger.Config `mapstructure:"logger"`
	App    *AppConf       `mapstructure:"app"`
}

type AppConf struct {
	Host     string          `mapstructure:"host"`
	Port     string          `mapstructure:"port"`
	Database *storage.Config `mapstructure:"database"`
}

func NewConfig(configFile string) (Config, error) {
	conf := Config{}
	_, err := config.ReadConfigFile(configFile, "toml", &conf)
	return conf, err
}
