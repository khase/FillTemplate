# Overview

`FillTemplate` is a little go tool which allowes you to simply fill Templates with one or more json/yml files

# Usage:
  FillTemplate [global flags] fromfile [flags] [values file path] 

## Flags:
  -f, --format string   input file format (supports json or yaml) (default "yaml")
  -h, --help            help for fromfile

## Global Flags:
| short command | long command | description |
|---------------|--------------|-------------|
|      | --exec string      | command to execute when output file changed (requires output filename to be set) |
|  -o, | --output string    | file to save the result to (default "stdout") |
|  -t, | --template string  | Template file to fill |

# Example

`FillTemplate -t [template file] -o [output file] fromfile -f json [input values file 0] ...`

When providing just one input file the content is available in the template via the root object `.`.
When providing multiple input files the root object will be an array of the different file contents and the separate contents can be extracted by the index function of the templating engine.
e.g (get the first files content ) 

`{{- $firstFile := index . 0 }}`

# Additional Info
Basic template usage: https://golangforall.com/en/post/templates.html

Additional Sprig functions: http://masterminds.github.io/sprig/
