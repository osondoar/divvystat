package models

type DivvyStatus struct {
	ExecutionTime   string
	StationBeanList map[int]Station
}
