package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type XMLPath struct {
	Fill string `xml:"fillColor,attr"`
	Data string `xml:"pathData,attr"`
}

type XMLData struct {
	XMLName xml.Name  `xml:"vector"`
	Width   float32   `xml:"viewportWidth,attr"`
	Height  float32   `xml:"viewportHeight,attr"`
	Paths   []XMLPath `xml:"path"`
}

type SVGPath struct {
	Fill string `xml:"fill,attr"`
	Data string `xml:"d,attr"`
}

type SVGData struct {
	XMLName xml.Name  `xml:"svg"`
	Viewbox string    `xml:"viewBox,attr"`
	Paths   []SVGPath `xml:"path"`
}

func main() {
	input := flag.String("in", "", "xml input")
	flag.Parse()

	raw, err := os.ReadFile(*input)
	if err != nil {
		log.Fatal(err)
	}

	parsed := XMLData{}

	buf := bytes.NewBuffer(raw)
	if err := xml.NewDecoder(buf).Decode(&parsed); err != nil {
		log.Fatal(err)
	}

	converted := SVGData{
		Viewbox: fmt.Sprintf("0 0 %v %v", parsed.Width, parsed.Height),
		Paths:   []SVGPath{},
	}

	for _, path := range parsed.Paths {
		converted.Paths = append(converted.Paths, SVGPath{
			Fill: "#444",
			Data: path.Data,
		})
	}

	var bufOut bytes.Buffer

	if err := xml.NewEncoder(&bufOut).Encode(&converted); err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(strings.ReplaceAll(*input, "xml", "svg"), bufOut.Bytes(), 0644); err != nil {
		log.Fatal(err)
	}
}
