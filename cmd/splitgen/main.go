package main

import (
	"encoding/xml"
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/jfoster/go-livesplit/lss"
)

func main() {
	flag.Parse()
	args := flag.Args()

	infile, err := os.OpenFile(args[0], os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer infile.Close()

	var splits lss.Run
	if err := xml.NewDecoder(infile).Decode(&splits); err != nil {
		log.Panic(err)
	}

	val, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatalf("cannot parse value %s to integer", args[1])
	}

	segments := make([]lss.Segment, val)
	for i := range segments {
		segments[i] = lss.Segment{
			Name: strconv.Itoa(i),
		}
	}

	splits.Segments = segments

	outfile, err := os.Create(args[0])
	if err != nil {
		log.Panic(err)
	}
	defer infile.Close()

	outfile.WriteString(xml.Header)

	if err := xml.NewEncoder(outfile).Encode(splits); err != nil {
		log.Panic(err)
	}
}
