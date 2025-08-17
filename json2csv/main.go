package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"os"
	"sort"
)

func main() {
	data := []map[string]string{}

	raw, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("read input: %v\n", err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(raw)).Decode(&data); err != nil {
		log.Fatalf("parse input as json: %v\n", err)
	}

	out := [][]string{}

	for i, obj := range data {
		if i == 0 {
			// the first row should be the keys
			row := []string{}
			for key := range obj {
				row = append(row, key)
			}
			sort.Strings(row)
			out = append(out, row)
		}

		row := []string{}

		for _, key := range out[0] {
			row = append(row, obj[key])
		}

		out = append(out, row)
	}

	if err := csv.NewWriter(os.Stdout).WriteAll(out); err != nil {
		log.Fatalf("write csv output: %v\n", err)
	}
}
