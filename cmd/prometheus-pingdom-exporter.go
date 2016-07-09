package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "prometheus-pingdom-exporter",
	Short: "prometheus-pingdom-exporter exports Pingdom metrics to Prometheus",
}
