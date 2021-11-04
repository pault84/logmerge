package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	in           []string
	out          string
	earliestDate time.Time = time.Now()
	pattern                = regexp.MustCompile(`\d{4}[-\/]\d{2}[-\/]\d{2}[T_ ]*[ ]?\d{1,}:\d{1,}:\d{1,}Z*,*[0-9]*`)
)

func mergeLogs(files []string) error {
	logrus.Infof("In: %v", files)

	openFiles := make([][]string, 3)

	for i := 0; i < len(files); i++ {
		f, err := os.Open(files[i])
		if err != nil {
			logrus.Fatalf("failed to open file %v with err: %v", files[i], err)
			return err
		}

		// The bufio.NewScanner() function is called in which the
		// object os.File passed as its parameter and this returns a
		// object bufio.Scanner which is further used on the
		// bufio.Scanner.Split() method.
		scanner := bufio.NewScanner(f)
		// The bufio.ScanLines is used as an
		// input to the method bufio.Scanner.Split()
		// and then the scanning forwards to each
		// new line using the bufio.Scanner.Scan()
		// method.

		scanner.Split(bufio.ScanLines)
		var text []string

		for scanner.Scan() {
			text = append(text, scanner.Text())
		}

		// The method os.File.Close() is called
		// on the os.File object to close the file
		f.Close()

		openFiles[i] = text
	}

	// Loop over files and purge files that do not match our regex.
	for i := 0; i < len(openFiles); i++ {
		for j, line := range openFiles[i] {
			matches := pattern.FindAllString(line, -1)
			if len(matches) == 0 {
				openFiles[i] = append(openFiles[i][:j], openFiles[i][j+1:]...)
			}
		}
	}

	// Recurse through all the files until we are told to stop.
	logLoop(openFiles)

	return nil
}

func logLoop(openFiles [][]string) {
	keepGoing := false
	earliestIndex := 999
	earliestDate = time.Now()

	// Loop over files
	for i := 0; i < len(openFiles); i++ {
		if len(openFiles[i]) == 0 {
			openFiles = append(openFiles[:i], openFiles[i+1:]...)
		} else {

			keepGoing = true
			line := openFiles[i][0]

			matches := pattern.FindAllString(line, -1)
			if len(matches) != 0 {

				timestamp := strings.Replace(matches[0], " ", "T", 1)
				timestamp = strings.ReplaceAll(timestamp, "/", "-")

				if strings.Contains(timestamp, "_") {
					if !strings.Contains(timestamp, "T") {
						timestamp = strings.Replace(timestamp, "_", "T", 1)
					} else {
						timestamp = strings.Replace(timestamp, "_", "", 1)
					}
				}

				// Change `2021-10-04 23:48:20,261` format into the expected `2021-10-04T23:48:03Z` format for parsing
				if strings.Contains(timestamp, ",") {
					splits := strings.Split(timestamp, ",")
					timestamp = splits[0]
				}

				// Add trailing Z if we do not have one to match the required format.
				if !strings.Contains(timestamp, "Z") {
					timestamp = timestamp + "Z"
				}

				// Parse the date into our date layout.
				t, err := time.Parse(time.RFC3339, timestamp)

				if err != nil {
					logrus.Infof("Line: %s", line)
					logrus.Fatalf("failed to parse date %s into a valid date: %v", timestamp, err)
				}

				// Is the parsed date after the earliest date we found?
				if t.Before(earliestDate) {
					earliestDate = t
					earliestIndex = i
				}
			} else {
				openFiles[i] = append(openFiles[i][:0], openFiles[i][1:]...)
			}
		}
	}

	if earliestIndex != 999 {
		fmt.Println(openFiles[earliestIndex][0])
		// Pop line from the array.
		openFiles[earliestIndex] = openFiles[earliestIndex][1:]
	}

	if keepGoing {
		logLoop(openFiles)
	}
}
