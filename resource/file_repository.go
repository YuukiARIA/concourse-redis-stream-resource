package resource

import (
	"io/ioutil"
	"path/filepath"
)

type FileRepository interface {
	Write(path string, content string) error
}

type fileSystemRepository struct {
	basePath string
}

func NewFileSystemRepository(basePath string) fileSystemRepository {
	return fileSystemRepository{basePath: basePath}
}

func (f fileSystemRepository) Write(path string, content string) error {
	outputPath := filepath.Join(f.basePath, path)
	return ioutil.WriteFile(outputPath, []byte(content), 0644)
}
