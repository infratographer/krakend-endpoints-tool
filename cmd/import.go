/*
Copyright © 2023 Infratographer Authors
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// importCmd represents the import command.
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "imports endpoints from a krakend configuration",
	Long: `Given a krakend configuration, this command will import the endpoints
from the configuration and generate the API definitions.

The "endpoints" parameter is the directory where the API definitions will be.`,
	RunE: importMain,
}

//nolint:gochecknoinits // this is a cobra command
func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.Flags().StringP("config", "c", "", "krakend configuration template")
}

func importMain(cmd *cobra.Command, args []string) error {
	endpoints := cmd.Flag("endpoints").Value.String()
	cfg := cmd.Flag("config").Value.String()

	return importEndpoint(endpoints, cfg)
}

func importEndpoint(endpoints, cfg string) error {
	if endpoints == "" {
		return fmt.Errorf("endpoints directory is required")
	}

	if cfg == "" {
		return fmt.Errorf("configuration template is required")
	}

	// Read the configuration file

	// Get the endpoints from the configuration

	// Iterate over the endpoints
	// e.g. using the WalkEndpoints function

	// Persist the endpoints to the endpoints directory

	return nil
}
