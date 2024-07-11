package filehandler

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type FileHandlerInterface interface {
	GetFilenames(sitePath string, dirName string) ([]string, error)
	ReadFile(filePath string) ([]byte, error)
}

type FileHandler struct {
}

// Ensure FileHandler implements FileHandlerInterface
var _ FileHandlerInterface = &FileHandler{}

func NewFileHandler() *FileHandler {
	return &FileHandler{}
}

func (fh FileHandler) GetFilenames(sitePath string, dirName string) ([]string, error) {
	var filenames []string
	dirPath := filepath.Join(sitePath, dirName)
	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			filenames = append(filenames, filepath.Base(path))
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %v", dirName, err)
	}
	return filenames, nil
}

func (fh FileHandler) ReadFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}
