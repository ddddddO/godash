package secretstore

import (
	"bufio"
	"os"
)

type fileSecretStore struct {
	path string
}

func NewFileSecretStore(path string) *fileSecretStore {
	return &fileSecretStore{
		path: path,
	}
}

func (s fileSecretStore) Load(typ string) (interface{}, error) {
	// typ?
	_ = typ

	f, err := os.Open(s.path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	secret := ""
	for scanner.Scan() {
		secret = scanner.Text()
	}

	return secret, nil
}
