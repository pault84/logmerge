package cmd

import (
	"bufio"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	in           []string
	out          string
	earliestDate time.Time = time.Now()
)

func mergeCommand() *cobra.Command {
	var mergeCmd = &cobra.Command{
		Use:   "merge",
		Short: "Merge the provided files",
		Long:  `Merge provided files based on timestamp and output new logfile.`,
		Run: func(cmd *cobra.Command, args []string) {
			mergeLogs(in, out)
		},
	}

	mergeCmd.Flags().StringSliceVarP(&in, "files", "f", []string{}, "comma-seperated logfiles that you wish to merge.")
	mergeCmd.Flags().StringVarP(&out, "output", "o", "merge.log", "outputlogfile.")

	return mergeCmd
}

func mergeLogs(files []string, output string) error {
	logrus.Infof("In: %v", in)
	logrus.Infof("Out: %v", out)

	openFiles := make([][]string, 3)
	mergedlog, err := os.OpenFile(output, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Fatalf("failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(mergedlog)
	defer mergedlog.Close()

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

	// Recurse through all the files until we are told to stop.
	logLoop(datawriter, openFiles)

	return nil
}

func logLoop(dataWriter *bufio.Writer, openFiles [][]string) {
	keepGoing := false
	earliestIndex := 999
	earliestDate = time.Now()

	// Loop over files
	for i := 0; i < len(openFiles); i++ {
		if len(openFiles[i]) > 0 && openFiles[i][0] == "" {
			openFiles = append(openFiles[:i], openFiles[i+1:]...)
		} else {
			if len(openFiles[i]) > 0 {
				keepGoing = true
				line := openFiles[i][0]

				regex := regexp.MustCompile(`\d{4}-\d{2}-\d{2}[T ]\d{2}:\d{2}:\d{2}Z*,*[0-9]*`)
				matches := regex.FindAllString(line, -1)
				if len(matches) != 0 {

					timestamp := strings.Replace(matches[0], " ", "T", 1)

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
	}

	if earliestIndex != 999 {
		_, err := dataWriter.WriteString(openFiles[earliestIndex][0] + "\n")
		if err != nil {
			logrus.Fatalf("failed to write line to file: %v", err)
		}

		// Pop line from the array.
		openFiles[earliestIndex] = openFiles[earliestIndex][1:]
		dataWriter.Flush()
	}

	if keepGoing {
		logLoop(dataWriter, openFiles)
	}
}
