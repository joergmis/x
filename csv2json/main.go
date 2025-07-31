package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"log"
	"os"
	"strings"
	"text/template"
)

func main() {
	infile := flag.String("in", "", "path to the csv file")
	flag.Parse()

	raw, err := os.ReadFile(*infile)
	if err != nil {
		log.Fatal(err)
	}

	reader := csv.NewReader(bytes.NewBuffer(raw))
	data, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	items := []map[string]string{}
	keys := data[0]

	for _, item := range data[1:] {
		object := map[string]string{}

		for i, key := range keys {
			object[key] = strings.ReplaceAll(item[i], "\"", "'")
			object[key] = strings.ReplaceAll(object[key], "\n", "")
		}

		items = append(items, object)
	}

	tpl, err := template.New("csv2json").Parse(tmpl)
	if err != nil {
		log.Fatal(err)
	}

	var out bytes.Buffer
	if err := tpl.Execute(&out, struct{ Items []map[string]string }{Items: items}); err != nil {
		log.Fatal(err)
	}

	rawJson := json.RawMessage(out.Bytes())
	b, err := json.MarshalIndent(&rawJson, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(strings.ReplaceAll(*infile, ".csv", ".json"), b, 0644); err != nil {
		log.Fatal(err)
	}
}

var tmpl = `[
	{{ $first := true }}
	{{ range .Items }}
	{{ if not $first }},{{ else }}{{ $first = false }}{{ end }}{
	{{ $firstItem := true }}
	{{ range $key, $value := . }}
	{{ if not $firstItem }},{{ else }}{{ $firstItem = false }}{{ end }}
	"{{ $key }}": "{{ $value }}"
	{{ end }}
	}
	{{ end }}
]`
