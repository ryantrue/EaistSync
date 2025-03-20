package files

import (
	"io/ioutil"
	"os"
)

// ReadFile возвращает содержимое файла по указанному пути.
func ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

// WriteFile записывает данные в файл по указанному пути с указанными правами доступа.
func WriteFile(path string, data []byte, perm os.FileMode) error {
	return ioutil.WriteFile(path, data, perm)
}
