package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"

	"github.com/jfoster/go-livesplit/lss"
)

const (
	layout = "01/02/2006 15:04:05" // DO NOT CHANGE
)

var (
	deleteDuration = flag.Duration("dd", time.Duration(time.Second*10), "delete attempts with less than this duration in seconds (default 0s)")
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		log.Fatalln(fmt.Errorf("Incorrect number of arguments. Expected: 1 Recieved: %d", flag.NArg()))
	}

	fmt.Println(*deleteDuration)

	infile, err := os.Open(args[0])
	if err != nil {
		log.Panic(err)
	}
	defer infile.Close()

	var splits lss.Run
	if err := xml.NewDecoder(infile).Decode(&splits); err != nil {
		log.Panic(err)
	}

	// sort attempt history by started time
	sort.SliceStable(splits.AttemptHistory, func(a, b int) bool {
		timeA := splits.AttemptHistory[a].Started.Time()
		timeB := splits.AttemptHistory[b].Started.Time()

		return timeA.Before(timeB)
	})

	// delete attempts that finish within 15 seconds of starting
	var attempts = splits.AttemptHistory
	for i, v := range attempts {
		if v.Started.Time().Add(*deleteDuration).After(v.Ended.Time()) {
			copy(attempts[i:], attempts[i+1:])    // Shift a[i+1:] left one index.
			attempts = attempts[:len(attempts)-1] // Truncate slice.
		}
	}
	splits.AttemptHistory = attempts

	for s, segment := range splits.Segments {
		var times = segment.Times

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

		// remove times that do not corrospond to an attempt
		for t := len(times) - 1; t >= 0; t-- {
			var found = false
			var time = times[t]
			for _, attempt := range splits.AttemptHistory {
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
		splits.Segments[s].Times = times
	}

	var attemptNum int64
	for a, attempt := range splits.AttemptHistory {
		attemptNum = int64(a) + 1
		for s, segment := range splits.Segments {
			for t, time := range segment.Times {
				if time.ID == attempt.ID {
					splits.Segments[s].Times[t].ID = attemptNum
				}
			}
		}
		splits.AttemptHistory[a].ID = attemptNum
	}
	splits.AttemptCount = attemptNum

	// file outputting

	filename := regexp.MustCompile("[\\\\/:*?\"<>|]").ReplaceAllString(fmt.Sprintf("%s - %s", splits.GameName, splits.CategoryName), "")

	data, err := xml.MarshalIndent(splits, "", "  ")
	if err != nil {
		log.Panic(err)
	}
	data = regexp.MustCompile(`(<[\w\s="]+)><(/)\w+(>)`).ReplaceAll(data, []byte("$1$2$3"))

	outfile, err := os.Create(filename + filepath.Ext(infile.Name()))
	if err != nil {
		log.Panic(err)
	}
	defer outfile.Close()

	outfile.WriteString(xml.Header)
	outfile.Write(data)
}
