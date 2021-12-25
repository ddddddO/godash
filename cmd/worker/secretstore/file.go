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

func (s file) Load(typ string) (interface{}, error) {
	f, err := os.Open(s.path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	secret := ""
	for scanner.Scan() {
		secret = scanner.Text()
		if strings.Contains(secret, typ) {
			break
		}
	}
	if secret == "" {
		return nil, errors.New("no datasource secret")
	}

	return secret, nil
}
