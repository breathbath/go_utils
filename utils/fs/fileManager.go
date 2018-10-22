package fs

import (
	"bufio"
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

func RTrimDirPath(inputDir string) string {
	return strings.TrimRight(inputDir, DS)
}

func ReadTextFile(filePath string, lineParser func(line string, lineNumber int) error) error {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return err
	}

	reader := bufio.NewReader(file)
	var line string
	lineNumber := -1
	for {
		lineNumber++
		line, err = reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSuffix(line, "\n")

		err = lineParser(line, lineNumber)
		if err != nil {
			return err
		}
	}

	if err != io.EOF {
		return err
	}

	return nil
}

func ReadFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return []byte{}, err
	}

	reader := bufio.NewReader(file)
	data, _, err := reader.ReadLine()
	return data, err
}

func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	return fileInfo.IsDir(), err
}
