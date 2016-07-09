package cmd

import (
	"github.com/spf13/cobra"
)

// RootCmd is the main Cobra Command.
var RootCmd = &cobra.Command{
	Use:   "prometheus-pingdom-exporter",
	Short: "prometheus-pingdom-exporter exports Pingdom metrics to Prometheus",
}
