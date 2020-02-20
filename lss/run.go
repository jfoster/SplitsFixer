package lss

// Run is a representation of xml data
type Run struct {
	Version      string `xml:"version,attr"`
	GameIcon     Icon
	GameName     string
	CategoryName string
	Metadata     struct {
		Run struct {
			ID string `xml:"id,attr"`
		} `xml:"Run"`
		Platform struct {
			Value        string  `xml:",chardata"`
			UsesEmulator lssBool `xml:"usesEmulator,attr"`
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
		Start          lssBool
		Reset          lssBool
		Split          lssBool
		CustomSettings []struct {
			ID    string `xml:"id,attr"`
			Type  string `xml:"type,attr"`
			Value string `xml:",chardata"`
		} `xml:"CustomSettings>Setting"`
	}
}

func (a Attempt) GetTimes(s []Segment) (times []*Time) {
	for _, v := range s {
		for _, t := range v.Times {
			if t.ID == a.ID {
				times = append(times, &t)
			}
		}
	}
	return
}

type Attempt struct {
	ID              int64   `xml:"id,attr"`
	Started         lssTime `xml:"started,attr"`
	IsStartedSynced lssBool `xml:"isStartedSynced,attr"`
	Ended           lssTime `xml:"ended,attr"`
	IsEndedSynced   lssBool `xml:"isEndedSynced,attr"`
	RealTime        string  `xml:"RealTime,omitempty"`
	GameTime        string  `xml:"GameTime,omitempty"`
}

type Segment struct {
	Name            string
	Icon            Icon        `xml:"Icon"`
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

func (t *Time) GetAttempt(attempts []Attempt) *Attempt {
	for _, v := range attempts {
		if v.ID == t.ID {
			return &v
		}
	}
	return nil
}

type Icon struct {
	Data string `xml:",cdata"`
}
