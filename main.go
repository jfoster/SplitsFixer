package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"time"
)

const (
	layout = "01/02/2006 15:04:05" // DO NOT CHANGE
)

func main() {
	lss := NewLSS()

	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal(fmt.Errorf("Incorrect number of arguments specied"))
	}

	data, err := ioutil.ReadFile(args[0])
	if err != nil {
		log.Fatal(err)
	}
	xml.Unmarshal(data, &lss.Run)

	sort.SliceStable(lss.Run.AttemptHistory, func(a, b int) bool {
		timeA, _ := time.Parse(layout, lss.Run.AttemptHistory[a].Started)
		timeB, _ := time.Parse(layout, lss.Run.AttemptHistory[b].Started)
		return timeA.Before(timeB)
	})

	var attemptHistory = lss.Run.AttemptHistory

	for s, segment := range lss.Run.Segments {
		var times = lss.Run.Segments[s].Times

		// sort times by ID value
		sort.SliceStable(times, func(a, b int) bool {
			return times[a].ID < times[b].ID
		})

		// remove ID values less than 1
		for i := len(times) - 1; i >= 0; i-- {
			if times[i].ID < 1 {
				times = append(times[:i], times[i+1:]...)
			}
		}

		// remove times that do not corrospond to a run
		for t := len(times) - 1; t >= 0; t-- {
			var found = false
			var time = times[t]
			for a, attempt := range attemptHistory {
				_ = a
				if time.ID == attempt.ID {
					found = true
					break
				}
			}
			if !found {
				times = append(times[:t], times[t+1:]...)
				fmt.Println(segment.Name, time.ID, "not found")
			}
		}

		lss.Run.Segments[s].Times = times
	}

	for a, attempt := range lss.Run.AttemptHistory {
		attemptNum := int64(a) + 1

		for s, segment := range lss.Run.Segments {
			for t, time := range segment.Times {
				if time.ID == attempt.ID {
					lss.Run.Segments[s].Times[t].ID = attemptNum
				}
			}
		}

		attemptHistory[a].ID = attemptNum
	}

	lss.Run.AttemptHistory = attemptHistory

	attemptCount := lss.Run.AttemptHistory[len(lss.Run.AttemptHistory)-1].ID
	lss.Run.AttemptCount = attemptCount
	fmt.Println(attemptCount)

	lss.WriteFile()
}
