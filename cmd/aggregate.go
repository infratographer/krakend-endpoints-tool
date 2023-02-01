/*
Copyright © 2023 Infratographer Authors
*/
package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
)

// aggregateCmd represents the aggregate command.
var aggregateCmd = &cobra.Command{
	Use:   "aggregate",
	Short: "Aggregates the API definitions into a single endpoints file",
	Long: `Aggregates the API definitions into a single endpoints file

This will take all of the API definitions and aggregate them into a single
file. This is useful for generating the endpoints configuration that will
eventually be mounted by the API Gateway.
`,
	RunE:         aggregateMain,
	SilenceUsage: true,
}

//nolint:gochecknoinits // this is a cobra command
func init() {
	rootCmd.AddCommand(aggregateCmd)

	aggregateCmd.Flags().StringP("output", "o", "-", "output file. defaults to stdout")
	aggregateCmd.Flags().BoolP("vhost", "v", false, "prepend the vhost to the endpoint")
}

func prependPrefix(obj any, prefix string, vhost bool) any {
	objMap, ok := obj.(map[string]any)
	if !ok {
		debug("not a map, skipping: %T", obj)
		return obj
	}
	suffixAny, ok := objMap["endpoint"]
	if !ok {
		debug("no endpoint to prepend: %v", objMap)
		return obj
	}

	suffix, ok := suffixAny.(string)
	if !ok {
		debug("endpoint is not a string: %v", suffixAny)
		return obj
	}

    endp := path.Join("/", prefix, suffix)

	if vhost {
		endp = path.Join("__virtual", endp)
	}

	objMap["endpoint"] = "/" + endp

	return objMap
}

func parseEndpoints(endpoints string, exceptions []string, vhost bool) ([]any, error) {
	endpts := []any{}
	err := WalkEndpoints(endpoints, exceptions, func(path string, typ endpointType, obj any, prefix string) error {
		switch typ {
		case arrayEndpoint:
			endptArr, ok := obj.([]any)
			if !ok {
				return fmt.Errorf("unexpected error: expected array of endpoints, got %T", obj)
			}

			for i, endpt := range endptArr {
				endptArr[i] = prependPrefix(endpt, prefix, vhost)
			}

			endpts = append(endpts, endptArr...)
		case objectEndpoint:
			obj := prependPrefix(obj, prefix, vhost)
			endpts = append(endpts, obj)
		default:
			return fmt.Errorf("unknown endpoint type: %s", path)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return endpts, nil
}

func aggregateMain(cmd *cobra.Command, args []string) error {
	endpoints := cmd.Flag("endpoints").Value.String()
	outf := cmd.Flag("output").Value.String()
	vhost := cmd.Flag("vhost").Value.String() == "true"
	return aggregate(endpoints, outf, vhost)
}

func aggregate(endpoints, outf string, vhost bool) error {
	if endpoints == "" {
		return fmt.Errorf("endpoints directory is required")
	}

	outfile, err := getOutputFile(outf)
	if err != nil {
		return err
	}

	defer outfile.Close()

	fmt.Println("# Aggregating endpoints in", green(endpoints))

	endpts, err := parseEndpoints(endpoints, exceptions, vhost)
	if err != nil {
		return err
	}

	if err := persistJSON(outfile, endpts); err != nil {
		return fmt.Errorf("failed to persist endpoints: %w", err)
	}

	// Print the endpoints to stdout if we're in debug mode
	if *debugMode && outf != "-" {
		info("# Aggregated endpoints:")
		if err := persistJSON(os.Stdout, endpts); err != nil {
			return fmt.Errorf("failed to persist endpoints: %w", err)
		}
	}

	return nil
}
