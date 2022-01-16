package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	inputFormat string

	fromfileCmd = &cobra.Command{
		Use:   "fromfile [values file path]",
		Short: "fromfile loads a values file and injects it into the templating engine",
		Long:  `fromfile loads a values file and injects it into the templating engine`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			templateString := getTemplate()
			inputValues := getInput(args[0])

			tmpl, err := template.New("Template").Parse(templateString)
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			err = tmpl.Execute(os.Stdout, inputValues)
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
		},
	}
)

func init() {
	fromfileCmd.Flags().StringVarP(&inputFormat, "format", "f", "yaml", "input file format (supports json or yaml)")
}

func getInput(inputFile string) interface{} {
	b, err := ioutil.ReadFile(inputFile)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	var result interface{}

	switch inputFormat {
	case "json":
		err = json.Unmarshal(b, &result)
	case "yaml":
		err = yaml.Unmarshal(b, &result)
	default:
		err = yaml.Unmarshal(b, &result)
	}

	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	return result
}
