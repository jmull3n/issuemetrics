package cmd

import (
	"github.com/jmull3n/issuemetrics/pkg/util"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	// for dump stack trace support
	_ "github.com/songgao/stacktraces/on/SIGUSR2"
)

// these values are set from the go build, do not change them
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// RootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "issuemetrics",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var versionCmd = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version, commit, date)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(v, c, d string) {
	version = v
	commit = c
	date = d
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
func getString(cmd *cobra.Command, key string) string {
	val, _ := cmd.Flags().GetString(key)
	return val
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.PersistentFlags().String("profile-port", util.Getenv("PROFILE_PORT", ""), "port to run profiler")
	rootCmd.PersistentFlags().String("metrics-port", util.Getenv("METRICS_PORT", "8083"), "port to run metrics")

}
