package fs

import (
	"bufio"
	"github.com/breathbath/go_utils/utils/errs"
	io2 "github.com/breathbath/go_utils/utils/io"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const DS = string(os.PathSeparator)

func GetCurrentPath() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir + string(os.PathSeparator), err
}

func ReadFilesInDirectory(dirPath string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dirPath)
}

func FileExistsOrFail(filePath string) {
	_, err := os.Stat(filePath)
	errs.FailOnError(err)
}

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func MkDir(dirPath string) error {
	if FileExists(dirPath) {
		return nil
	}

	return os.MkdirAll(dirPath, os.ModePerm)
}

func RmFile(filePath string) bool {
	if !FileExists(filePath) {
		return true
	}

	err := os.Remove(filePath)
	return err == nil
}

func JoinPath(parts ...string) string {
	return filepath.Join(parts...)
}

func RTrimDirPath(inputDir string) string {
	return strings.TrimRight(inputDir, DS)
}

func ReadTextFile(filePath string, lineParser func(line string, lineNumber int) error) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	lineNumber := -1
	for {
		lineNumber++
		line, readErr := reader.ReadString('\n')
		if line == "" && readErr != nil {
			break
		}

		line = strings.TrimSuffix(line, "\n")

		err = lineParser(line, lineNumber)
		if err != nil {
			return err
		}

		if readErr != nil {
			err = readErr
			break
		}
	}

	if err == io.EOF {
		return nil
	}

	return err
}

func ReadFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	defer func() {
		errClose := file.Close()
		if errClose != nil {
			io2.OutputError(errClose, "", "")
		}
	}()

	if err != nil {
		return []byte{}, err
	}

	reader := bufio.NewReader(file)
	data, _, err := reader.ReadLine()
	return data, err
}

func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), nil
}
