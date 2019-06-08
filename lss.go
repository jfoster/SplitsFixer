package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
)

type LSS struct {
	Run Run
}

func NewLSS() LSS {
	lss := LSS{}
	return lss
}

func (lss LSS) WriteFile() {
	data, err := xml.MarshalIndent(lss.Run, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// prepend xml header to data
	data = append([]byte(xml.Header), data...)

	reg, err := regexp.Compile("[\\\\/:*?\"<>|]")
	if err != nil {
		log.Fatal(err)
	}

	fileName := fmt.Sprintf("%s - %s.lss", lss.Run.GameName, lss.Run.CategoryName)

	err = ioutil.WriteFile(reg.ReplaceAllString(fileName, ""), data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
