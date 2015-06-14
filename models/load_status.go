package models

type LoadStatus struct {
	Time string `json:"execution_time"`
	Load int    `json:"load"`
}
