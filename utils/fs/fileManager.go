package fs

import (
	"bufio"
	"fmt"
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

func TouchFile(fullFilePath string) error {
	err := ioutil.WriteFile(fullFilePath, []byte{}, 0644)
	return err
}

func WriteFileString(fullFilePath, data string, perm os.FileMode) error {
	err := ioutil.WriteFile(fullFilePath, []byte(data), 0644)
	return err
}

func ReadFileString(fullFilePath string) (string, error) {
	data, err := ioutil.ReadFile(fullFilePath)
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
	defer func() {
		err := source.Close()
		if err != nil {
			io2.OutputWarning("", "failed to close source file %s: %v", fullFilePath, err)
		}
	}()

	destination, err := os.Create(destFilePath)
	if err != nil {
		return fmt.Errorf("cannot open dest file '%s': %v", destFilePath, err)
	}
	defer func() {
		err := destination.Close()
		if err != nil {
			io2.OutputWarning("", "failed to close dest file %s: %v", destFilePath, err)
		}
	}()

	_, err = io.Copy(destination, source)
	if err != nil {
		return fmt.Errorf("cannot copy '%s' to '%s': %v", fullFilePath, destFilePath, err)
	}

	return nil
}
