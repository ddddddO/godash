package model

type Task struct {
	Kind           kind   `json:"kind"`
	Settings       string `json:"settings"`
	DataSourceType string `json:"type"`
	Query          string `json:"query"`
}

// Taskの種類
type kind uint

const (
	KindSettings kind = iota
	KindQuery
)
