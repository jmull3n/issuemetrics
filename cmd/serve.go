package cmd

import (
	"syscall"
	"os/signal"
	"context"
	"os"

	"github.com/jmull3n/issuemetrics/pkg/metrics"
	"github.com/jmull3n/issuemetrics/pkg/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "fire up the issuemetrics api",
	Run: func(cmd *cobra.Command, args []string) {
		// setup logging
		log.SetOutput(os.Stderr)
		log.SetFormatter(&log.JSONFormatter{})
		log.SetLevel(log.InfoLevel)

		// setup context
		ctx := context.Background()

		// parse cli args
		port := getString(cmd, "port")
		metricsPort := getString(cmd, "metrics-port")

		// start the real server
		server.Start(ctx, port)

		// start the metric server
		metrics.StartServer(ctx, metricsPort)

		end := gracefulShutdown()
		<-end

	},
}

func gracefulShutdown() <-chan struct{} {
	end := make(chan struct{})
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-s
		log.Info("Sutting down gracefully.")
		close(end)
	}()
	return end
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().String("port", "8000", "the port which exposes the rest api")
}
