package cmd

import (
  "github.com/spf13/cobra"
  "github.com/sirupsen/logrus"
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
		  logrus.Info("Hugo Static Site Generator v0.9 -- HEAD")
		},
	}

	mergeCmd.Flags().StringSliceVar(&in, "f", []string{}, "comma-seperated logfiles that you wish to merge.")
	mergeCmd.Flags().StringVar(&out, "o", "", "outputlogfile.")

	return mergeCmd
}

func mergeLogs(files []string, output string) error {
	openFiles := []string

	mergedlog, err := os.OpenFile("test.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Fatalf("failed creating file: %s", err)
	}

	//datawriter := bufio.NewWriter(mergedlog)
	defer mergedlog.Close()

	for i := 0; i < len(files); i++ {
		f, err := os.Open(files[i])
		if err != nil {
			logrus.Fatalf("failed to open file %v with err: %v", files[i], err)
		}

		// The bufio.NewScanner() function is called in which the
		// object os.File passed as its parameter and this returns a
		// object bufio.Scanner which is further used on the
		// bufio.Scanner.Split() method.
		scanner = bufio.NewScanner(f)
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

	while len(openFiles) > 0 {
		//earliestIndex := 999
		//earliestDate := time.Now()

		// Loop over files
		for i := 0; i < len(openFiles); i++ {
			line := openfiles[i][0]
			logrus.Infof("Line: %s", line)
			//Get time from line.
			
		}

		
	}
 /*
	// Write line
	_, err = datawriter.WriteString(openfiles[earliestIndex] + "\n")
	if err != nil {
		logrus.Fatalf("failed to open file %v with err: %v", files[i], err)
	}
	// Pop line from the array.
	openfiles[earliestIndex] = openfiles[earliestIndex][1:]

	// Flush data to disk
	datawriter.Flush()
*/
	return nil
}