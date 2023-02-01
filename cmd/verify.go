/*
Copyright © 2023 Infratographer Authors
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xeipuuv/gojsonschema"
)

// verifyCmd represents the verify command.
var verifyCmd = &cobra.Command{
	Use:          "verify",
	Short:        "Verifies API endpoint configurations",
	Long:         `Verifies API endpoint configurations`,
	RunE:         verifyMain,
	SilenceUsage: true,
}

//nolint:gochecknoinits // this is a cobra command
func init() {
	rootCmd.AddCommand(verifyCmd)
}

func validateObj(path string, obj any, endptldr gojsonschema.JSONLoader) error {
	// load json file
	endptfile := gojsonschema.NewGoLoader(obj)
	if endptfile == nil {
		return fmt.Errorf("invalid endpoint configuration")
	}

	// validate json file
	result, err := gojsonschema.Validate(endptldr, endptfile)
	if err != nil {
		return err
	}

	if !result.Valid() {
		info("The document %s is not valid. see errors:", red(path))
		for _, desc := range result.Errors() {
			info(" - %s", red(desc))
		}

		return fmt.Errorf("invalid endpoint configuration: %s", path)
	}
	return nil
}

func validateArray(path string, arr []any, endptldr gojsonschema.JSONLoader) error {
	for _, obj := range arr {
		if err := validateObj(path, obj, endptldr); err != nil {
			return err
		}
	}
	return nil
}

func verifyMain(cmd *cobra.Command, args []string) error {
	endpoints := cmd.Flag("endpoints").Value.String()
	return verify(endpoints)
}

func verify(endpoints string) error {
	if endpoints == "" {
		return fmt.Errorf("endpoints directory is required")
	}

	// load json schema
	endptldr := gojsonschema.NewReferenceLoader("https://www.krakend.io/schema/endpoint.json")
	if endptldr == nil {
		return fmt.Errorf("invalid endpoint schema")
	}

	err := WalkEndpoints(endpoints, []string{}, func(path string, typ endpointType, obj any, _ string) error {
		switch typ {
		case arrayEndpoint:
			objAny, ok := obj.([]any)
			if !ok {
				return fmt.Errorf("invalid endpoint configuration: %s", path)
			}
			return validateArray(path, objAny, endptldr)
		case objectEndpoint:
			return validateObj(path, obj, endptldr)
		default:
			return fmt.Errorf("unknown endpoint type: %s", path)
		}
	})
	if err != nil {
		return err
	}

	info(green("* All endpoints are valid"))
	return nil
}
