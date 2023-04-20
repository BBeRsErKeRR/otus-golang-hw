package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/app"
	versioncmd "github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/cmd"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/scheduler"
	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Short: "Scheduler application",
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

		storage := app.GetEventUseCase(config.Database)
		err = storage.Connect(ctx)
		if err != nil {
			logg.Error("Error create db connection: " + err.Error())
			return
		}

		produce := scheduler.GetProducer(config.Producer, logg)
		appl := scheduler.New(logg, scheduler.GetProducerUseCase(produce), storage, config.Duration)

		if err = produce.Connect(ctx); err != nil {
			logg.Error("Error create mq connection: " + err.Error())
		}

		go func() {
			if err := appl.Run(ctx); err != nil {
				logg.Error("failed to consume mq: " + err.Error())
				cancel()
			}
		}()

		defer produce.Close(ctx)
		defer storage.Close(ctx)

		<-ctx.Done()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "./configs/scheduler_config.toml", "Configuration file path")
	rootCmd.AddCommand(versioncmd.VersionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
