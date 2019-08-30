package fs

import (
	"bufio"
	"encoding/csv"
	"fmt"
	io2 "github.com/breathbath/go_utils/utils/io"
	"io"
	"os"
)

func ReadCsvFile(filePath string, lineParser func(int, []string) error, lineSep rune, leastColumnsCount int) (err error) {

	if !FileExists(filePath) {
		err = fmt.Errorf("No file found at %s", filePath)
		return
	}

	csvFile, err := os.Open(filePath)
	if err != nil {
		return
	}

	defer func() {
		closeErr := csvFile.Close()
		if closeErr != nil {
			io2.OutputError(closeErr, "", "")
		}
	}()

	fileReader := csv.NewReader(bufio.NewReader(csvFile))
	fileReader.Comma = lineSep

	lineNumber := 0

	for {
		lineNumber++
		line, curErr := fileReader.Read()
		if curErr == io.EOF {
			break
		} else if curErr != nil {
			return fmt.Errorf(
				"Failed to parse csv file %s at line %d, reason: %s",
				filePath,
				lineNumber,
				curErr.Error(),
			)
		}

		if len(line) != leastColumnsCount {
			return fmt.Errorf(
				"Unexpected columns count %d on line %d: expected count is %d",
				len(line),
				lineNumber,
				leastColumnsCount,
			)
		}

		err = lineParser(lineNumber, line)
		if err != nil {
			return
		}
	}

	return
}
