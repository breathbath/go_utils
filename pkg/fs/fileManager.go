package fs

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	io2 "github.com/breathbath/go_utils/v3/pkg/io"

	"github.com/breathbath/go_utils/v3/pkg/errs"
)

const DS = string(os.PathSeparator)
const defaultPerm = os.FileMode(0600)

func GetCurrentPath() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir + string(os.PathSeparator), err
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

func TouchFile(fullFilePath string) error {
	err := os.WriteFile(fullFilePath, []byte{}, defaultPerm)
	return err
}

func WriteFileString(fullFilePath, data string, perm os.FileMode) error {
	err := os.WriteFile(fullFilePath, []byte(data), defaultPerm)
	return err
}

func ReadFileString(fullFilePath string) (string, error) {
	data, err := os.ReadFile(fullFilePath)
	return string(data), err
}

func ReadFileStringSecure(fullFilePath string) string {
	data, _ := ReadFileString(fullFilePath)
	return data
}

func CopyFile(fullFilePath, destFilePath string) error {
	source, err := os.Open(fullFilePath)
	if err != nil {
		return fmt.Errorf("cannot open file '%s': %v", fullFilePath, err)
	}
	defer io2.CloseResourceSecure("source file", source)

	destination, err := os.Create(destFilePath)
	if err != nil {
		return fmt.Errorf("cannot open dest file '%s': %v", destFilePath, err)
	}
	defer io2.CloseResourceSecure("dest file", destination)

	_, err = io.Copy(destination, source)
	if err != nil {
		return fmt.Errorf("cannot copy '%s' to '%s': %v", fullFilePath, destFilePath, err)
	}

	return nil
}
