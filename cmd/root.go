/*
Copyright © 2023 Infratographer Authors
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "krakend-endpoints-tool",
	Short: "krakend-endpoints-tool is a simple tool to help lint and aggregate the APIs in this repo",
	Long: `krakend-endpoints-tool is a simple tool to help lint and aggregate the APIs in this repo

It's meant to be used in CI to ensure that the APIs are consistent and that the
documentation is up to date.
`,
}

//nolint:gochecknoinits // this is a cobra command
func init() {
	rootCmd.PersistentFlags().String("endpoints", "", "endpoints directory")
	debugMode = rootCmd.PersistentFlags().Bool("debug", false, "debug mode")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
