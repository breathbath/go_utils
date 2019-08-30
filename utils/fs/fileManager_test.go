package fs

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestGetCurrentPath(t *testing.T) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	assert.NoError(t, err)
	expectedCurPathString := dir + string(os.PathSeparator)

	curPathString, err := GetCurrentPath()
	assert.NoError(t, err)

	assert.Equal(t, expectedCurPathString, curPathString)
}

func TestReadFilesInDirectory(t *testing.T) {
	err := MkDir("lala")
	assert.NoError(t, err)

	err = ioutil.WriteFile("lala/someFile.txt", []byte(""), 0644)
	assert.NoError(t, err)

	files, err := ReadFilesInDirectory("lala")
	assert.NoError(t, err)

	assert.Len(t, files, 1)
	for _, file := range files {
		assert.Equal(t, "someFile.txt", file.Name())
	}

	err = os.RemoveAll("lala")
	assert.NoError(t, err)
}

func TestFileExistsOrFail(t *testing.T) {
	assert.PanicsWithValue(t, "stat someNonExistingFile.txt: no such file or directory", func() {
		FileExistsOrFail("someNonExistingFile.txt")
	})
}

func TestMkDir(t *testing.T) {
	err := os.Mkdir("lala_exists", 0644)
	assert.NoError(t, err)

	err = MkDir("lala_exists")
	assert.NoError(t, err)

	err = MkDir("lala")
	assert.NoError(t, err)

	dStat, err := os.Stat("lala")
	assert.NoError(t, err)

	assert.True(t, dStat.IsDir())

	err = os.RemoveAll("lala")
	assert.NoError(t, err)

	err = os.RemoveAll("lala_exists")
	assert.NoError(t, err)
}

func TestRmFile(t *testing.T) {
	err := ioutil.WriteFile("someFileToDelete.txt", []byte{}, 0644)
	assert.NoError(t, err)

	result := RmFile("someFileToDelete.txt")
	assert.True(t, result)

	if !result {
		err = os.Remove("someFileToDelete.txt")
		assert.NoError(t, err)
	}

	result = RmFile("nonExistingFile.txt")
	assert.True(t, result)
}

func TestJoinPath(t *testing.T) {
	result := JoinPath("lsls", "mama")
	assert.Equal(t, "lsls/mama", result)
}

func TestRTrimDirPath(t *testing.T) {
	result := RTrimDirPath("lsls/")
	assert.Equal(t, "lsls", result)
}

func TestReadTextFile(t *testing.T) {
	testFile(t, "testFile.txt", []byte("line1\nline2\nline3"), func(t *testing.T) {
		allLines := map[int]string{}
		err := ReadTextFile(
			"testFile.txt",
			func(line string, lineNumb int) error {
				allLines[lineNumb] = line
				return nil
			},
		)
		assert.NoError(t, err)

		expectedLines := map[int]string{
			0: "line1",
			1: "line2",
			2: "line3",
		}
		assert.EqualValues(t, expectedLines, allLines)
	})

	testFile(t, "testFile.txt", []byte("line1\n"), func(t *testing.T) {
		allLines := map[int]string{}
		err := ReadTextFile(
			"testFile.txt",
			func(line string, lineNumb int) error {
				allLines[lineNumb] = line
				return nil
			},
		)
		assert.NoError(t, err)

		expectedLines := map[int]string{
			0: "line1",
		}
		assert.EqualValues(t, expectedLines, allLines)
	})

	testFile(t, "testFileWithErr.txt", []byte("11"), func(t *testing.T) {
		err := ReadTextFile(
			"testFileWithErr.txt",
			func(line string, lineNumb int) error {
				return errors.New("Some error")
			},
		)
		assert.EqualError(t, err, "Some error")
	})

	err := ReadTextFile(
		"nonExistingFile.txt",
		func(line string, lineNumb int) error {
			return nil
		},
	)

	assert.Error(t, err)
}

func TestReadFile(t *testing.T) {
	testFile(t, "testFile2.txt", []byte("someText"), func(t *testing.T) {
		data, err := ReadFile("testFile2.txt")
		assert.NoError(t, err)
		assert.Equal(t, "someText", string(data))
	})

	data, err := ReadFile("nonExistingFile")
	assert.Error(t, err)
	assert.Equal(t, []byte{}, data)
}

func testFile(t *testing.T, fileName string, fileContent []byte, testFunc func(t *testing.T)) {
	RmFile(fileName)

	err := ioutil.WriteFile(fileName, fileContent, 0644)
	assert.NoError(t, err)

	testFunc(t)

	err = os.Remove(fileName)
	assert.NoError(t, err)
}

func TestIsDirectory(t *testing.T) {
	err := MkDir("lala2")
	assert.NoError(t, err)

	isDir, err := IsDirectory("lala2")
	assert.NoError(t, err)
	assert.True(t, isDir)

	err = os.RemoveAll("lala2")
	assert.NoError(t, err)

	isDir, err = IsDirectory("nonExistingDir")
	assert.Error(t, err)
	assert.False(t, isDir)

	testFile(t, "someNonDirFile", []byte(""), func(t *testing.T) {
		isDir, err = IsDirectory("someNonDirFile")
		assert.NoError(t, err)
		assert.False(t, isDir)
	})
}