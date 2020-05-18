package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/den-is/ktempl/pkg/logging"
	"github.com/den-is/ktempl/pkg/validation"
	"github.com/den-is/ktempl/pkg/worker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Initializes job, listens for signals
func Gates(cmd *cobra.Command, args []string) {

	// TODO: add central place for validation logic
	// Check if template file exists
	if err := validation.CheckFileExists(viper.GetString("template")); err != nil {
		logging.LogWithFields(
			logging.Fields{
				"component": "starter",
			}, "error", "Template file at given path does not exist.", err)
		os.Exit(1)
	}

	signalCh := make(chan os.Signal, 1)

	go worker.Worker()

	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	for s := range signalCh {
		if s == syscall.SIGINT || s == syscall.SIGTERM {
			logging.LogWithFields(
				logging.Fields{
					"component": "starter",
				}, "info", fmt.Sprintf("Got %q signal: exiting", s))
			os.Exit(0)
		}
	}

}