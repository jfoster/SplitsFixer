package lss

import (
	"encoding/xml"
	"fmt"
	"regexp"

	"github.com/jfoster/file"
)

var (
	origFile file.File
)

// LSS type
type LSS struct {
	Run Run
}

// NewLSSFile creates an LSS object from a file
func New(filename string) (lss LSS, err error) {
	err = lss.ReadFile(file.New(filename))
	return lss, err
}

// ReadFile reads from a file into LSS object
func (lss *LSS) ReadFile(f file.File) error {
	origFile = f

	data, err := f.Read()
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

	reg := regexp.MustCompile("[\\\\/:*?\"<>|]")
	newFileName := reg.ReplaceAllString(fmt.Sprintf("%s - %s", lss.Run.GameName, lss.Run.CategoryName), "")

	f := file.File{Dir: origFile.Dir, Name: newFileName, Extension: origFile.Extension}

	err = f.Write(data, 0644)
	if err != nil {
		return err
	}

	return nil
}

type CDATA struct {
	Cdata string `xml:",cdata"`
}

// Run is a representation of xml data
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
