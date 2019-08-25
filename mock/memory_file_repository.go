package mock

type memoryFileRepository struct {
	contents map[string]string
}

func NewMemoryFileRepository() memoryFileRepository {
	return memoryFileRepository{contents: map[string]string{}}
}

func (m memoryFileRepository) Write(path, content string) error {
	m.contents[path] = content
	return nil
}

func (m memoryFileRepository) Count() int {
	return len(m.contents)
}

func (m memoryFileRepository) GetContent(path string) string {
	return m.contents[path]
}
