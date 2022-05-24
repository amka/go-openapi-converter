package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "openapi2conv <input.json> <output.yaml>",
		Short: "Gopher CLI in Go",
		Long:  `Gopher CLI application written in Go.`,
		Args:  cobra.MinimumNArgs(2),
		// Uncomment the following line if your bare application
		// has an action associated with it:
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Print: " + strings.Join(args, " "))
			convert(args[0], args[1])
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func convert(input_json, output_yaml string) {
	input, err := os.ReadFile(input_json)
	if err != nil {
		panic(err)
	}

	var doc openapi2.T
	if err = json.Unmarshal(input, &doc); err != nil {
		panic(err)
	}

	outputJSON, err := json.Marshal(doc)
	if err != nil {
		panic(err)
	}
	var docAgainFromJSON openapi2.T
	if err = json.Unmarshal(outputJSON, &docAgainFromJSON); err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(doc, docAgainFromJSON) {
		fmt.Println("objects doc & docAgainFromJSON should be the same")
	}

	outputYAML, err := yaml.Marshal(doc)
	if err != nil {
		panic(err)
	}
	var docAgainFromYAML openapi2.T
	if err = yaml.Unmarshal(outputYAML, &docAgainFromYAML); err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(doc, docAgainFromYAML) {
		fmt.Println("objects doc & docAgainFromYAML should be the same")
	}

	// Write outputYAML to file
	if err := os.WriteFile(output_yaml, outputYAML, 0644); err != nil {
		panic(err)
	}

	fmt.Print("Successfully converted " + input_json + " to " + output_yaml)
}
