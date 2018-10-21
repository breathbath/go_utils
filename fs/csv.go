package fs

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func ReadCsvFile(filePath string, lineParser func(int, []string) error, lineSep rune, leastColumnsCount int) (err error) {
	if !FileExists(filePath) {
		err = fmt.Errorf("No file found at %s", filePath)
		return
	}

	csvFile, err := os.Open(filePath)
	defer csvFile.Close()
	if err != nil {
		return
	}

	fileReader := csv.NewReader(bufio.NewReader(csvFile))
	fileReader.Comma = lineSep

	lineNumber := 0
	for {
		lineNumber++
		line, curErr := fileReader.Read()
		if curErr == io.EOF {
			break
		} else if curErr != nil {
			err = fmt.Errorf("Failed to parse csv file %s at line %d, reason: %s", filePath, lineNumber, curErr.Error())
			return
		}

		if len(line) < leastColumnsCount {
			return fmt.Errorf("Unexpected columns number %d: it should be at least %d", len(line), leastColumnsCount)
		}

		err = lineParser(lineNumber, line)
		if err != nil {
			return
		}
	}

	return
}

