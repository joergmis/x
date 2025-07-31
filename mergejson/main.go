package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"os"
	"strings"
)

/*
 notes
 - works (only?) with basic arrays with map-ish objects
 - the order of the input files is important - only the items of the first file are extended
 - the matcher is case sensitive
*/

func main() {
	in := flag.String("in", "", "list of files to merge, split by comma")
	matcher := flag.String("matcher", "", "string on which the data should be compared and combined")
	flag.Parse()

	files := strings.Split(*in, ",")
	if len(files) < 2 {
		log.Fatal("not enough files to merge")
	}
	if *matcher == "" {
		log.Fatal("empty matcher")
	}

	maps := [][]map[string]string{}

	for _, file := range files {
		data := []map[string]string{}

		raw, err := os.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}

		if err := json.NewDecoder(bytes.NewBuffer(raw)).Decode(&data); err != nil {
			log.Fatal(err)
		}

		maps = append(maps, data)
	}

	merged := []map[string]string{}

	first := maps[0]
	for _, data := range maps[1:] {
		for _, value := range first {
			for _, val := range data {
				if val[*matcher] == value[*matcher] {
					combined := map[string]string{}

					for key, val := range val {
						combined[key] = val
					}

					for key, val := range value {
						combined[key] = val
					}

					merged = append(merged, combined)

				}
			}
		}
	}

	b, err := json.MarshalIndent(&merged, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(strings.ReplaceAll(files[0], ".json", ".merged.json"), b, 0644); err != nil {
		log.Fatal(err)
	}
}
