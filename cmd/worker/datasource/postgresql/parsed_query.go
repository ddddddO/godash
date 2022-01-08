package postgresql

import (
	"strings"

	"github.com/pkg/errors"
)

type queryType uint

const (
	undefined queryType = iota
	selectType
	insertType
	updateType
	deleteType
)

var (
	errUndefinedType = errors.New("sql query is undefined dml")
)

type parsedQuery struct {
	qType queryType
	value string
}

func newParsedQuery(rawQuery string) *parsedQuery {
	return &parsedQuery{
		value: strings.TrimSpace(rawQuery),
	}
}

func (pq *parsedQuery) validate() error {
	switch {
	case strings.HasPrefix(pq.value, "select"):
		return nil
	case strings.HasPrefix(pq.value, "insert"):
		return nil
	case strings.HasPrefix(pq.value, "update"):
		return nil
	case strings.HasPrefix(pq.value, "delete"):
		return nil
	}
	return errUndefinedType
}

func (pq *parsedQuery) decideQueryType() {
	qType := undefined
	switch {
	case strings.HasPrefix(pq.value, "select"):
		qType = selectType
	case strings.HasPrefix(pq.value, "insert"):
		qType = insertType
	case strings.HasPrefix(pq.value, "update"):
		qType = updateType
	case strings.HasPrefix(pq.value, "delete"):
		qType = deleteType
	}
	pq.qType = qType
}
