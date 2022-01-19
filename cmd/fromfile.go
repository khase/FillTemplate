package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Masterminds/sprig"
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
			var inputValuesList []interface{}
			for _, content := range args {
				inputValuesList = append(inputValuesList, getInput(content))
			}

			var input interface{}
			if len(inputValuesList) == 1 {
				input = inputValuesList[0]
			} else {
				input = inputValuesList
			}

			tmpl, err := template.New("Template").Funcs(sprig.FuncMap()).Parse(templateString)
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}

			buf := new(bytes.Buffer)

			err = tmpl.Execute(buf, input)
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}

			switch outputFile {
			case "":
				log.Println(buf.String())
			case "stdout":
				log.Println(buf.String())
			default:
				changed, err := dumpToFile(outputFile, buf.String())
				if err != nil {
					log.Fatalln(err)
					os.Exit(1)
				}
				if changed && executeOnChangeCommand != "" {
					log.Println("Output changed, executing post command")
					cmd := exec.Command(executeOnChangeCommand)
					var out bytes.Buffer
					cmd.Stdout = &out

					err := cmd.Run()
					if err != nil {
						log.Fatalln(err)
						os.Exit(1)
					}
					log.Printf("Command returned:\n%s\n", string(out.String()))
				}
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
func dumpToFile(name string, data string) (cahnged bool, err error) {
	if err != nil {
		return false, err
	}

	if _, err := os.Stat(name); errors.Is(err, os.ErrNotExist) {
		dir := filepath.Dir(name)
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return false, err
		}

		err = writeFile(name, data)
		if err != nil {
			return false, err
		}

		return true, nil
	}

	currentContent, err := readFile(name)
	if err != nil {
		return false, err
	}

	if currentContent == data {
		return false, nil
	}

	err = writeFile(name, data)
	if err != nil {
		return false, err
	}

	return true, nil
}

func writeFile(name string, content string) error {
	f, err := os.Create(name)

	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(content)

	if err != nil {
		return err
	}

	return nil
}

func readFile(name string) (string, error) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
