package main

import (
	"encoding/xml"
	"flag"
	"log"
	"os"

	"github.com/jfoster/go-livesplit/lss"

	"github.com/go-yaml/yaml"
)

type Note struct {
	Name string   `yaml:"name"`
	Text []string `yaml:"text"`
}

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		log.Fatalln("Incorrect number of args")
	}

	infile, err := os.Open(args[0])
	if err != nil {
		log.Panic(err)
	}
	defer infile.Close()

	var splits lss.Run
	if err := xml.NewDecoder(infile).Decode(&splits); err != nil {
		log.Panic(err)
	}

	var notes = make([]Note, len(splits.Segments))

	for i, v := range splits.Segments {
		notes[i] = Note{
			Name: v.Name,
			Text: []string{"sample text"},
		}
	}

	if err := yaml.NewEncoder(os.Stdout).Encode(notes); err != nil {
		log.Panic(err)
	}
}
