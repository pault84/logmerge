package cmd

import (
	"bufio"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	in  []string
	out string
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

	iteration := 0
	for {
		filesWithLogs := 0
		earliestIndex := 999
		earliestDate := time.Now()

		// Loop over files
		for i := 0; i < len(openFiles); i++ {
			if len(openFiles[i]) > 0 && openFiles[i][0] == "" {
				openFiles = append(openFiles[:i], openFiles[i+1:]...)
			} else {
				if len(openFiles[i]) > 0 {
					filesWithLogs++
					line := openFiles[i][0]
					if strings.Contains(line, "time=") {
						cols := strings.Split(line, " ")
						for j, col := range cols {
							if strings.Contains(col, "time=") {
								// Take the 2nd index because the first one will be empty.
								timestamp := strings.Split(col, "time=")[1]
								if !strings.Contains(timestamp, ":") {
									if strings.Contains(cols[j+1], ":") {
										timestamp = timestamp + "T" + cols[j+1]
									}
								}
								// Trim the quotes from the string.
								timestamp = strings.Trim(timestamp, "\"")
								// Parse the date into our date layout.
								t, err := time.Parse(time.RFC3339, timestamp)
								if err != nil {
									logrus.Infof("Line: %s", line)
									logrus.Fatalf("failed to parse date %s into a valid date: %v", timestamp, err)
								}
								// Is the parsed date after the earliest date we found?
								if earliestDate.After(t) {
									earliestDate = t
									earliestIndex = i
								}
							}
						}
					} else {
						openFiles[i] = append(openFiles[i][:0], openFiles[i][1:]...)
					}
				}
			}
		}

		iteration++
		if filesWithLogs == 0 {
			break
		}

		if earliestIndex != 999 {
			// Write line
			_, err = datawriter.WriteString(openFiles[earliestIndex][0] + "\n")
			if err != nil {
				logrus.Fatalf("failed to write line from file %v with err: %v", files[earliestIndex], err)
			}

			// Pop line from the array.
			openFiles[earliestIndex] = openFiles[earliestIndex][1:]

			datawriter.Flush()
		}
	}
	return nil
}
