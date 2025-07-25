package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var (
		defaultDirectory        = "."
		defaultPathExcludes     = ""
		defaultPathIncludes     = ""
		defaultTemplateLocation = "./copyright.tpl"

		templateLocation = flag.String("tpl", defaultTemplateLocation, "Path to the copyright header template")
		pathExcludes     = flag.String("exclusions", defaultPathExcludes, "List of strings (to search in the file path) to exclude from the check, separated by commas")
		pathIncludes     = flag.String("inclusions", defaultPathIncludes, "List of strings (to search in the file path) to include from the check, separated by commas")
		directory        = flag.String("dir", defaultDirectory, "Directory to check the files for missing copyright headers")

		fix = flag.Bool("fix", false, "Automatically prepend the copyright in missing files")

		exclusions = []string{}
		inclusions = []string{}
	)

	flag.Parse()

	raw, err := os.ReadFile(*templateLocation)
	if err != nil {
		log.Fatal(err)
	}
	copyrightText := strings.TrimSpace(string(raw))

	missingCopyrights := []string{}

	exclusions = strings.Split(*pathExcludes, ",")
	inclusions = strings.Split(*pathIncludes, ",")

	filepath.WalkDir(*directory, func(path string, d os.DirEntry, e error) error {
		for _, exclusion := range exclusions {
			if strings.Contains(path, exclusion) {
				return nil
			}
		}

		includeMatches := false
		for _, inclusion := range inclusions {
			if strings.Contains(path, inclusion) {
				includeMatches = true
			}
		}

		if !includeMatches {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		if !strings.HasPrefix(string(content), copyrightText) {
			missingCopyrights = append(missingCopyrights, path)
		}

		return nil
	})

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range missingCopyrights {
		path := filepath.Join(currentDir, file)

		if *fix {
			content, err := os.ReadFile(path)
			if err != nil {
				log.Fatal(err)
			}

			if err := os.WriteFile(path, []byte(fmt.Sprintf("%s\n%s", string(copyrightText), string(content))), 0644); err != nil {
				log.Fatal(err)
			}

			log.Printf("fixed copyright: %v\n", path)
		} else {
			log.Printf("missing copyright: %v\n", path)
		}
	}

	if len(missingCopyrights) != 0 {
		os.Exit(-1)
	}
}
