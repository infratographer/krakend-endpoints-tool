/*
Copyright © 2023 Infratographer Authors
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// generateCmd represents the generate command.
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generates a krakend configuration",
	Long: `Generates a krakend.io configuration.

Given a path to the endpoint definitions and a krakend configuration template,
this command will generate a krakend configuration file.

Given that krakend itself uses golang templates, we resort to setting up an
explicit identifier for replacing the endpoints in the template.`,
	RunE:         generateMain,
	SilenceUsage: true,
}

//nolint:gochecknoinits // this is a cobra command
func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringP("config", "c", "", "krakend configuration template")
	generateCmd.Flags().StringP("output", "o", "-", "output file. defaults to stdout")
	generateCmd.Flags().StringP("identifier", "i", "$ENDPOINTS$", "identifier for the endpoints in the template")
	generateCmd.Flags().BoolP("vhost", "v", false, "prepend the vhost to the endpoint")
}

func generateMain(cmd *cobra.Command, args []string) error {
	endpoints := cmd.Flag("endpoints").Value.String()
	cfg := cmd.Flag("config").Value.String()
	outf := cmd.Flag("output").Value.String()
	id := cmd.Flag("identifier").Value.String()
	vhost := cmd.Flag("vhost").Value.String() == "true"

	return Generate(endpoints, cfg, outf, id, vhost)
}

func Generate(endpoints, cfg, outf, id string, vhost bool) error {
	if endpoints == "" {
		return fmt.Errorf("endpoints directory is required")
	}

	if cfg == "" {
		return fmt.Errorf("configuration template is required")
	}

	outfile, err := getOutputFile(outf)
	if err != nil {
		return err
	}

	cfgbytes, err := os.ReadFile(cfg)
	if err != nil {
		return err
	}

	defer outfile.Close()

	endpts, err := parseEndpoints(endpoints, exceptions, vhost)
	if err != nil {
		return err
	}

	stringBuffer := &strings.Builder{}
	if err := persistJSON(stringBuffer, endpts); err != nil {
		return fmt.Errorf("error persisting the endpoints: %w", err)
	}

	replacer := strings.NewReplacer(id, strings.TrimSpace(stringBuffer.String()))
	cfgfull := strings.TrimSpace(string(cfgbytes))
	if _, err := replacer.WriteString(outfile, cfgfull); err != nil {
		return fmt.Errorf("error writing the configuration: %w", err)
	}

	if *debugMode && outf != "-" {
		if _, err := replacer.WriteString(os.Stdout, cfgfull); err != nil {
			return fmt.Errorf("error writing the configuration: %w", err)
		}
	}

	return nil
}
