package cmd

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	templateFile string
	outputFile   string

	rootCmd = &cobra.Command{
		Use:   "FillTemplate",
		Short: "FillTemplate is a go tool for simple template fill actions from the cli.",
		Long:  `FillTemplate is a go tool for simple template fill actions from the cli.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&templateFile, "template", "t", "", "Template file to fill")
	rootCmd.PersistentFlags().StringVarP(&outputFile, "output", "o", "stdout", "file to save the result to")

	rootCmd.MarkFlagRequired("template")

	rootCmd.AddCommand(fromfileCmd)
}

func getTemplate() string {
	if templateFile == "" {
		log.Fatal("No template file defined")
		os.Exit(1)
	}

	b, err := ioutil.ReadFile(templateFile)
	if err != nil {
		log.Print(err)
	}

	return string(b)
}
