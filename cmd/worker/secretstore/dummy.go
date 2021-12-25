package secretstore

import (
	"fmt"
)

type dummySecretStore struct {
}

func NewDummySecretStore() *dummySecretStore {
	return &dummySecretStore{}
}

func (dummySecretStore) Load(typ string) (interface{}, error) {
	fmt.Println("not yet impl")

	_ = typ
	return "dummy secret", nil
}