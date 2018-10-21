package fs

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestReadNonExistingCsvFile(t *testing.T) {
	err := ReadCsvFile(
		"non_existing.csv",
		func(i int, strings []string) error {
			return nil
		},
		',',
		1,
	)

	assert.EqualError(t, err, "No file found at non_existing.csv")
}

func TestReadExistingCsvFile(t *testing.T) {
	csvFileRaw := []byte("col1,col2,col3\n1,2,3\n4,5,6\n7,8,9")
	err := ioutil.WriteFile("testFile.csv", csvFileRaw, 0644)
	assert.NoError(t, err)

	allLines := map[int][]string{}
	err = ReadCsvFile(
		"testFile.csv",
		func(i int, strings []string) error {
			allLines[i]= strings
			return nil
		},
		',',
		3,
	)
	assert.NoError(t, err)

	expectedLines := map[int][]string{
		1: {"col1", "col2", "col3"},
		2: {"1", "2", "3"},
		3: {"4", "5", "6"},
		4: {"7", "8", "9"},
	}
	assert.EqualValues(t, expectedLines, allLines)

	os.Remove("testFile.csv")
}
