package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
)

func main() {
	matcher := flag.String("matcher", "", "string on which the data should be compared and combined")
	flag.Parse()

	data := []map[string]string{}

	raw, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("read input: %v\n", err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(raw)).Decode(&data); err != nil {
		log.Fatalf("parse input as json: %v\n", err)
	}

	cleaned := map[string]map[string]string{}

	for _, obj := range data {
		key := obj[*matcher]

		_, ok := cleaned[key]
		if !ok {
			cleaned[key] = obj
		} else {
			for k, value := range obj {
				val, ok := cleaned[key][k]
				if !ok || val == "" {
					cleaned[key][k] = value
				}
			}
		}
	}

	out := []map[string]string{}
	for _, value := range cleaned {
		out = append(out, value)
	}

	b, err := json.MarshalIndent(&out, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stdout.Write(b); err != nil {
		log.Fatal(err)
	}

}
