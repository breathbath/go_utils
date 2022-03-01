package fs

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCurrentPath(t *testing.T) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	assert.NoError(t, err)
	expectedCurPathString := dir + string(os.PathSeparator)

	curPathString, err := GetCurrentPath()
	assert.NoError(t, err)

	assert.Equal(t, expectedCurPathString, curPathString)
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
	err := TouchFile("someFileToDelete.txt")
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
				return errors.New("some error")
			},
		)
		assert.EqualError(t, err, "some error")
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

	err := os.WriteFile(fileName, fileContent, 0600)
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

func TestTouchFile(t *testing.T) {
	filePath := "fileToTouch.txt"
	assert.False(t, FileExists(filePath))

	err := TouchFile(filePath)
	assert.NoError(t, err)

	defer RmFile(filePath)

	assert.True(t, FileExists(filePath))

	data, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.Len(t, data, 0)
}

func TestReadWriteFileString(t *testing.T) {
	filePath := "fileToWrite2.txt"
	err := WriteFileString(filePath, "some", 0664)
	assert.NoError(t, err)

	defer RmFile(filePath)
	data, err := ReadFileString(filePath)
	assert.NoError(t, err)
	if err == nil {
		assert.Equal(t, "some", data)
	}

	data, err = ReadFileString("someUnknownFile.txt")
	assert.IsType(t, &os.PathError{}, err)
	assert.Equal(t, "", data)
}

func TestReadWriteFileStringSecure(t *testing.T) {
	data := ReadFileStringSecure("someUnknownFile.txt")
	assert.Equal(t, "", data)

	filePath := "fileToWrite.txt"
	defer RmFile(filePath)

	err := WriteFileString(filePath, "some", 0664)
	assert.NoError(t, err)

	if err != nil {
		return
	}

	actualData := ReadFileStringSecure(filePath)
	assert.Equal(t, "some", actualData)
}

func TestCopyFile(t *testing.T) {
	filePathSource := "sourceFile.txt"
	filePathDest := "destFile.txt"

	err := WriteFileString(filePathSource, "data to copy", 0664)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	defer RmFile(filePathSource)

	err = CopyFile(filePathSource, filePathDest)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	defer RmFile(filePathDest)

	actualData := ReadFileStringSecure(filePathDest)
	assert.Equal(t, "data to copy", actualData)

	err = CopyFile("nonExistingFile.txt", filePathDest)
	assert.EqualError(t, err, "cannot open file 'nonExistingFile.txt': open nonExistingFile.txt: no such file or directory")

	err = CopyFile(filePathSource, "/")
	assert.EqualError(t, err, "cannot open dest file '/': open /: is a directory")
}
