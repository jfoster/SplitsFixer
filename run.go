package main

type CDATA struct {
	CDATA string `xml:",cdata"`
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
