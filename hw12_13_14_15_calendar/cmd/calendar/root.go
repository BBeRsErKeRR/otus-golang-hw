package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/BBeRsErKeRR/otus-golang-hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Short: "Calendar application",
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

		storage := app.GetEventUseCase(config.App.Database)
		calendar := app.New(logg, storage)

		server := internalhttp.NewServer(logg, calendar, config.App.HttpServer)

		go func() {
			if err := server.Start(ctx); err != nil {
				logg.Error("failed to start http server: " + err.Error())
				cancel()
			}
		}()

		defer server.Stop()

		<-ctx.Done()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "./configs/config.toml", "Configuration file path")
	rootCmd.AddCommand(versionCmd)
}
