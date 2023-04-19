package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	version_cmd "github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/cmd"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/sender"
	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Short: "Sender application",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		defer cancel()

		config, err := NewConfig(cfgFile)
		if err != nil {
			log.Println("Error create config: " + err.Error())
			return
		}

		logg, err := logger.New(config.Logger)
		if err != nil {
			log.Println("Error create app logger: " + err.Error())
			return
		}

		consumer := sender.GetConsumer(config.Consumer, logg)
		err = consumer.Connect(ctx)
		if err != nil {
			logg.Error("Error create mq connection: " + err.Error())
			return
		}

		app := sender.New(logg, sender.GetConsumerUseCase(consumer))

		go func() {
			if err := app.Consume(ctx); err != nil {
				logg.Error("failed to consume mq: " + err.Error())
				cancel()
			}
		}()

		defer consumer.Close(ctx)

		<-ctx.Done()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "./configs/sender_config.toml", "Configuration file path")
	rootCmd.AddCommand(version_cmd.VersionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
