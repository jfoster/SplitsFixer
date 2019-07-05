package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
)

var (
	origFile string
)

// LSS type
type LSS struct {
	Run Run
}

// NewLSS creates an empty LSS object
func NewLSS() LSS {
	return LSS{}
}

// NewLSSFile creates an LSS object from a file
func NewLSSFile(filename string) (LSS, error) {
	lss := NewLSS()
	err := lss.ReadFile(filename)
	return lss, err
}

// ReadFile reads from a file into LSS object
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

// WriteFile writes a LSS object to a file
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

type CDATA struct {
	Cdata string `xml:",cdata"`
}

type Run struct {
	Version      string `xml:"version,attr"`
	GameIcon     CDATA
	GameName     string
	CategoryName string
	Metadata     struct {
		Run struct {
			ID string `xml:"id,attr"`
		} `xml:"Run"`
		Platform struct {
			Value        string `xml:",chardata"`
			UsesEmulator bool   `xml:"usesEmulator,attr"`
		} `xml:"Platform"`
		Region struct {
		} `xml:"Region"`
		Variables struct {
		} `xml:"Variables"`
	}
	Offset               string
	AttemptCount         int64
	AttemptHistory       []Attempt `xml:"AttemptHistory>Attempt"`
	Segments             []Segment `xml:"Segments>Segment"`
	AutoSplitterSettings struct {
		Version        string
		ScriptPath     string
		Start          bool
		Reset          bool
		Split          bool
		CustomSettings []struct {
			ID    string `xml:"id,attr"`
			Type  string `xml:"type,attr"`
			Value string `xml:",chardata"`
		} `xml:"CustomSettings>Setting"`
	}
}

type Attempt struct {
	ID              int64  `xml:"id,attr"`
	Started         string `xml:"started,attr"`
	IsStartedSynced bool   `xml:"isStartedSynced,attr"`
	Ended           string `xml:"ended,attr"`
	IsEndedSynced   bool   `xml:"isEndedSynced,attr"`
	RealTime        string `xml:"RealTime,omitempty"`
	GameTime        string `xml:"GameTime,omitempty"`
}

type Segment struct {
	Name            string
	Icon            CDATA
	SplitTimes      []SplitTime `xml:"SplitTimes>SplitTime"`
	BestSegmentTime struct {
		RealTime string `xml:"RealTime,omitempty"`
		GameTime string `xml:"GameTime,omitempty"`
	}
	Times []Time `xml:"SegmentHistory>Time"`
}
type SplitTime struct {
	Name     string `xml:"name,attr"`
	RealTime string `xml:"RealTime,omitempty"`
	GameTime string `xml:"GameTime,omitempty"`
}

type Time struct {
	ID       int64  `xml:"id,attr"`
	RealTime string `xml:"RealTime,omitempty"`
	GameTime string `xml:"GameTime,omitempty"`
}
