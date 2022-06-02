package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
)

func Execute() {
	var outputType string

	var rootCmd = &cobra.Command{
		Use:   "go-openapi-converter input_file [output_file]",
		Short: "OpenApi 2 to 3",
		Long:  `Convert OpenApi 2 to OpenApi 3`,
		Args:  cobra.MinimumNArgs(1),
		// Uncomment the following line if your bare application
		// has an action associated with it:
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				convert(args[0], "", outputType)
			} else {
				convert(args[0], args[1], outputType)
			}
		},
	}
	rootCmd.PersistentFlags().StringVar(&outputType, "type", "json", "Output file type (yaml|json)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Convert converts a JSON file to a YAML file.
// `inputPath` is the path to the JSON file.
// `outputPath` is the path to the YAML file. If no `outputPath` is specified,
// the output is printed to stdout.
// `outputType` is the output file type, it could be `json` or `yaml`.
func convert(inputPath, outputPath, outputType string) {

	// Resolve filename
	inputFilepath, err := filepath.Abs(inputPath)
	if err != nil {
		panic(err)
	}

	input, err := os.ReadFile(inputFilepath)
	if err != nil {
		panic(err)
	}

	var doc openapi2.T
	if err = json.Unmarshal(input, &doc); err != nil {
		panic(err)
	}

	convertedDoc, err := openapi2conv.ToV3(&doc)

	// Validate the document
	outputJSON, err := json.Marshal(convertedDoc)
	if err != nil {
		panic(err)
	}
	var docAgainFromJSON openapi2.T
	if err = json.Unmarshal(outputJSON, &docAgainFromJSON); err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(convertedDoc, docAgainFromJSON) {
		fmt.Println("objects doc & docAgainFromJSON should be the same")
	}

	var outputBytes []byte

	// Detect what kind of file we are converting to
	if outputType == "yaml" {
		outputBytes, err = yaml.Marshal(convertedDoc)
		if err != nil {
			panic(err)
		}
	} else {
		// outputType = "json"
		outputBytes, err = json.Marshal(convertedDoc)
		if err != nil {
			panic(err)
		}
	}

	if len(outputPath) == 0 {
		fmt.Print(string(outputBytes))
	} else {
		// Write outputBytes to file
		if !strings.HasSuffix(outputPath, "."+outputType) {
			outputPath = outputPath + "." + outputType
		}

		if err := os.WriteFile(outputPath, outputBytes, 0644); err != nil {
			panic(err)
		}
	}
}
