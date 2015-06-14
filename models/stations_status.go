package models

type StationsStatus struct {
	ExecutionTime string
	stations      map[int]Station
}

func NewStationsStatus(executionTime string) StationsStatus {
	stationsStatus := StationsStatus{ExecutionTime: executionTime}
	stationsStatus.stations = make(map[int]Station)
	return stationsStatus
}

func (ds StationsStatus) Station(stationId int) (Station, bool) {
	station, ok := ds.stations[stationId]
	return station, ok
}

func (ds StationsStatus) AddStation(station Station) {
	ds.stations[station.Id] = station
}

func (ds StationsStatus) GetStations() map[int]Station {
	return ds.stations
}
