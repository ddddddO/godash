package model

type Task struct {
	DataSourceType string `json:"type"`
	Query          string `json:"query"`
}
