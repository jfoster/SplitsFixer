package lss

import "time"

// LSS type
type LSS struct {
	Run Run
}

type lssTime string

const customTimeFormat = `01/02/2006 15:04:05`

func (t lssTime) Time() time.Time {
	time, _ := time.Parse(customTimeFormat, string(t))
	return time
}

type lssBool string

func (b lssBool) Bool() bool {
	return b == "True"
}
