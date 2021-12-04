package model

type Result struct {
	StatusCode  int    `json:"status_code"`
	QueryResult string `json:"result"`
}
