package postgresql

import (
	"strings"

	"github.com/pkg/errors"
)

var (
	errUndefinedType = errors.New("sql query is undefined dml")
)

type parsedQuery struct {
	value string
}

func newParsedQuery(rawQuery string) *parsedQuery {
	return &parsedQuery{
		value: strings.TrimSpace(rawQuery),
	}
}

func (pq *parsedQuery) validate() error {
	switch {
	case pq.isSelect():
		return nil
	case pq.isInsert():
		return nil
	case pq.isUpdate():
		return nil
	case pq.isDelete():
		return nil
	}
	return errUndefinedType
}

func (pq *parsedQuery) isSelect() bool {
	return strings.HasPrefix(pq.value, "select")
}

func (pq *parsedQuery) isInsert() bool {
	return strings.HasPrefix(pq.value, "insert")
}

func (pq *parsedQuery) isUpdate() bool {
	return strings.HasPrefix(pq.value, "update")
}

func (pq *parsedQuery) isDelete() bool {
	return strings.HasPrefix(pq.value, "delete")
}
