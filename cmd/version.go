package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run:   versionRun,
	}

	version   string
	goVersion string
	gitCommit string
	osArch    string
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

func versionRun(cmd *cobra.Command, args []string) {
	fmt.Printf("Version:\t%v\n", version)
	fmt.Printf("Go version:\t%v\n", goVersion)
	fmt.Printf("Git commit:\t%v\n", gitCommit)
	fmt.Printf("OS/Arch:\t%v\n", osArch)
}
