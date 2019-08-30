package fs

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"os/exec"
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

func TestReadCsvFileWithBadPermissions(t *testing.T) {
	err := ioutil.WriteFile("testFileBadPermission.csv", []byte(""), 0644)
	assert.NoError(t, err)

	cmd := exec.Command("chmod", "000", "testFileBadPermission.csv")
	out, err := cmd.Output()
	fmt.Println(out)
	assert.NoError(t, err)

	err = ReadCsvFile(
		"testFileBadPermission.csv",
		func(i int, strings []string) error {
			return nil
		},
		',',
		1,
	)

	assert.EqualError(t, err, "open testFileBadPermission.csv: permission denied")

	err = os.Remove("testFileBadPermission.csv")
	assert.NoError(t, err)
}

func TestReadExistingCsvFile(t *testing.T) {
	csvFileRaw := []byte("col1,col2,col3\n1,2,3\n4,5,6\n7,8,9")
	err := ioutil.WriteFile("testFile.csv", csvFileRaw, 0644)
	assert.NoError(t, err)

	allLines := map[int][]string{}
	err = ReadCsvFile(
		"testFile.csv",
		func(i int, strings []string) error {
			allLines[i] = strings
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

	err = os.Remove("testFile.csv")
	assert.NoError(t, err)
}

func TestReadFaultyCsvFile(t *testing.T) {
	testWrongCsvFile(
		t,
		"col1,col2,col3\n1",
		"Failed to parse csv file testFaultyFile.csv at line 2, reason: record on line 2: wrong number of fields",
		nil,
		3,
	)

	testWrongCsvFile(
		t,
		"col1,col2,col3\n1,2,3",
		"Unexpected columns count 3 on line 1: expected count is 1",
		nil,
		1,
	)

	testWrongCsvFile(
		t,
		"col1,col2,col3\n1,2,3",
		"Some error",
		errors.New("Some error"),
		3,
	)
}

func testWrongCsvFile(t *testing.T, data, expectedError string, lineError error, expectedColsCount int) {
	t.Helper()

	err := ioutil.WriteFile("testFaultyFile.csv", []byte(data), 0644)
	assert.NoError(t, err)

	allLines := map[int][]string{}
	err = ReadCsvFile(
		"testFaultyFile.csv",
		func(i int, strings []string) error {
			allLines[i] = strings
			return lineError
		},
		',',
		expectedColsCount,
	)
	assert.EqualError(
		t,
		err,
		expectedError,
	)

	err = os.Remove("testFaultyFile.csv")
	assert.NoError(t, err)
}
