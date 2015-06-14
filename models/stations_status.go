package models

type StationsStatus struct {
	ExecutionTime      string
	ExecutionTimeEpoch int64
	stations           map[int]Station
}

func NewStationsStatus(executionTime string) StationsStatus {
	stationsStatus := StationsStatus{ExecutionTime: executionTime}
	stationsStatus.stations = make(map[int]Station)
	return stationsStatus
}

// func NewStationsStatusFromApi(apiStatus StationsStatusApi) StationsStatus {
// 	stationsStatus := StationsStatus{ExecutionTime: apiStatus.ExecutionTime}
// 	for _, station := range apiStatus.StationBeanList {
// 		station := models.Station{Id: station.Id, AvailableDocks: station.AvailableDocks}
// 		stationsStatus.AddStation(station)
// 	}
//
// 	stationsStatus.ExecutionTimeEpoch = getEpoch()
// 	return stationsStatus
// }

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
