package secretstore

import (
	"bufio"
	"os"
	"strings"

	"github.com/pkg/errors"
)

type file struct {
	path string
}

func NewFile(path string) *file {
	return &file{
		path: path,
	}
}

func (fi *file) Store(typ, settings string) error {
	f, err := os.Create(fi.path)
	if err != nil {
		return err
	}
	defer f.Close()

	// 既に保存されていればreturn
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), typ) {
			return nil
		}
	}

	if _, err := f.Write([]byte(settings)); err != nil {
		return err
	}
	return nil
}

func (fi *file) Load(typ string) (interface{}, error) {
	f, err := os.Open(fi.path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	secret := ""
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), typ) {
			secret = scanner.Text()
			break
		}
	}
	if secret == "" {
		return nil, errors.New("no datasource secret")
	}

	return secret, nil
}
