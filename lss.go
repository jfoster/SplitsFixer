package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"regexp"
)

type LSS struct {
	Run Run
}

func NewLSS() LSS {
	return LSS{}
}

func (lss LSS) ReadFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = xml.Unmarshal(data, &lss.Run)
	if err != nil {
		return err
	}

	return nil
}

func (lss LSS) WriteFile() error {
	data, err := xml.MarshalIndent(lss.Run, "", "  ")
	if err != nil {
		return err
	}

	// prepend xml header to data
	data = append([]byte(xml.Header), data...)

	reg, err := regexp.Compile("[\\\\/:*?\"<>|]")
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s - %s.lss", lss.Run.GameName, lss.Run.CategoryName)

	err = ioutil.WriteFile(reg.ReplaceAllString(fileName, ""), data, 0644)
	if err != nil {
		return err
	}

	return nil
}
