package services

import (
	"log"
	"strconv"
	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/osondoar/divvystat/models"
)

type LoadReporter struct {
	redisWrapper *RedisWrapper
}

func NewLoadReporter() *LoadReporter {
	var loadReporter LoadReporter
	loadReporter.redisWrapper = NewRedisWrapper()
	return &loadReporter
}

func (dr LoadReporter) GetAverageLoad(lastN int) (int, string) {
	divvyStatuses := dr.getLastStatuses(lastN)
	statusesSize := len(divvyStatuses)
	var average int

	for i := 0; i < statusesSize-1; i++ {
		diff := getStatusDiff(divvyStatuses[i], divvyStatuses[i+1])
		average += diff
	}

	if statusesSize > 1 {
		return (average / (statusesSize - 1)), divvyStatuses[0].ExecutionTime
	} else if statusesSize > 0 {
		return 0, divvyStatuses[0].ExecutionTime
	} else {
		return 0, ""
	}
}

func (dr LoadReporter) GetAverageLoads(from int64, to int64) []models.LoadStatus { //map[string]int {
	loadStatuses := []models.LoadStatus{}
	encodedLoads, err1 := redis.Strings(dr.redisWrapper.redisConn.Do("ZREVRANGEBYSCORE", loadKey, to, from))

	if err1 != nil {
		log.Print("Error: ", err1)
	}

	for _, encodedLoad := range encodedLoads {
		parts := strings.Split(encodedLoad, "__")
		load, _ := strconv.Atoi(parts[1])
		loadStatuses = append(loadStatuses, models.LoadStatus{Time: parts[0], Load: load})
	}

	return loadStatuses
}

func (dr LoadReporter) getLastStatuses(n int) []models.StationsStatus {
	var lastStatuses []models.StationsStatus

	lastStatusKeys, err := redis.Strings(dr.redisWrapper.redisConn.Do("ZREVRANGE", statusKeys, 0, n))
	if err != nil {
		log.Print("Error: ", err)
	}

	for _, statusKey := range lastStatusKeys {
		stationsStatus := models.NewStationsStatus(statusKey)

		statusTuples, err1 := redis.Ints(dr.redisWrapper.redisConn.Do("HGETALL", statusKey))
		if err1 != nil {
			log.Print("Error: ", err1)
		}

		for i := 0; i < len(statusTuples); i += 2 {
			station := models.Station{Id: statusTuples[i], AvailableDocks: statusTuples[i+1]}
			stationsStatus.AddStation(station)
		}
		lastStatuses = append(lastStatuses, stationsStatus)
	}

	return lastStatuses
}

func getStatusDiff(dr1 models.StationsStatus, dr2 models.StationsStatus) int {
	var diffAcum int
	for stationId, station1 := range dr1.GetStations() {
		s2, ok := dr2.Station(stationId)
		if !ok {
			continue
		}
		// Available docks
		diff := s2.AvailableDocks - station1.AvailableDocks
		if diff < 0 {
			diff = -diff
		}

		diffAcum += diff
	}

	return diffAcum
}
