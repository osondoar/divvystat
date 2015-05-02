package services

import (
	"log"
	"strconv"

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

	var average int
	for i := 0; i < lastN; i++ {
		if i >= len(divvyStatuses)-1 {
			break
		}
		diff := getStatusDiff(divvyStatuses[i], divvyStatuses[i+1])
		average += diff
	}

	return (average / lastN), divvyStatuses[0].ExecutionTime
}

func (dr LoadReporter) GetAverageLoads() map[string]int {
	loads := make(map[string]int)
	loadsTuples, err := redis.Strings(dr.redisWrapper.redisConn.Do("HGETALL", loadKeys))
	if err != nil {
		log.Print("Error: ", err)
	}

	for i := 0; i < len(loadsTuples); i += 2 {
		load, err := strconv.Atoi(loadsTuples[i+1])
		if err != nil {
			log.Fatal(err)
		}
		loads[loadsTuples[i]] = load
	}

	return loads
}

func (dr LoadReporter) getLastStatuses(n int) []models.DivvyStatus {
	var divvyStatuses []models.DivvyStatus

	lastStatuKeys, err := redis.Strings(dr.redisWrapper.redisConn.Do("ZREVRANGE", statusKeys, 0, n))
	if err != nil {
		log.Print("Error: ", err)
	}

	for _, statusKey := range lastStatuKeys {
		stations := make(map[int]models.Station)

		statusTuples, err := redis.Ints(dr.redisWrapper.redisConn.Do("HGETALL", statusKey))
		if err != nil {
			log.Print("Error: ", err)
		}

		for i := 0; i < len(statusTuples); i += 2 {
			stations[statusTuples[i]] = models.Station{Id: statusTuples[i], AvailableDocks: statusTuples[i+1]}
		}
		divvyStatuses = append(divvyStatuses, models.DivvyStatus{StationBeanList: stations, ExecutionTime: statusKey})
	}

	return divvyStatuses
}

func getStatusDiff(dr1 models.DivvyStatus, dr2 models.DivvyStatus) int {
	var diffAcum int
	for key, station1 := range dr1.StationBeanList {
		s2, ok := dr2.StationBeanList[key]
		if !ok {
			continue
		}
		diff := s2.AvailableDocks - station1.AvailableDocks
		if diff < 0 {
			diff = -diff
		}

		diffAcum += diff
	}

	return diffAcum
}
