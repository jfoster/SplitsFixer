package main

import (
	"flag"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/jfoster/splitfix/lss"
)

const (
	layout = "01/02/2006 15:04:05" // DO NOT CHANGE
)

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatalln(fmt.Errorf("Incorrect number of arguments. Expected: 1 Recieved: %d", flag.NArg()))
	}

	splits, err := lss.NewLSSFile(flag.Args()[0])
	if err != nil {
		log.Fatal(err)
	}

	sort.SliceStable(splits.Run.AttemptHistory, func(a, b int) bool {
		timeA, _ := time.Parse(layout, splits.Run.AttemptHistory[a].Started)
		timeB, _ := time.Parse(layout, splits.Run.AttemptHistory[b].Started)
		return timeA.Before(timeB)
	})

	var attemptHistory = splits.Run.AttemptHistory

	for s, segment := range splits.Run.Segments {
		var times = splits.Run.Segments[s].Times

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

		splits.Run.Segments[s].Times = times
	}

	for a, attempt := range splits.Run.AttemptHistory {
		attemptNum := int64(a) + 1

		for s, segment := range splits.Run.Segments {
			for t, time := range segment.Times {
				if time.ID == attempt.ID {
					splits.Run.Segments[s].Times[t].ID = attemptNum
				}
			}
		}

		attemptHistory[a].ID = attemptNum
	}

	splits.Run.AttemptCount = attemptHistory[len(attemptHistory)-1].ID
	splits.Run.AttemptHistory = attemptHistory

	err = splits.WriteFile()
	if err != nil {
		log.Fatal(err)
	}
}
