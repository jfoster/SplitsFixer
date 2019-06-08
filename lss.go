package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
)

type LSS struct {
	Run Run
}

func NewLSS() LSS {
	return LSS{}
}

func NewLSSFromFile(filename string) (LSS, error) {
	lss := NewLSS()
	err := lss.ReadFile(filename)
	return lss, err
}

var origFile string

func (lss *LSS) ReadFile(filename string) error {
	origFile = filename

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

func (lss *LSS) WriteFile() error {
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

	fileName := reg.ReplaceAllString(fmt.Sprintf("%s - %s.lss", lss.Run.GameName, lss.Run.CategoryName), "")
	filePath := filepath.Join(filepath.Dir(origFile), fileName)

	err = ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
