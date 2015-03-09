package services

import (
	"log"
	"time"
)

var divvyApiUrl = "http://www.divvybikes.com/stations/json"
var statusKeys = "statusKeys"
var loadKeys = "loads"

type LoadUpdater struct {
	redisHost      string
	redisWrapper   *RedisWrapper
	divvyApiClient DivvyApiClient
	loadReporter   *LoadReporter
}

func NewLoadUpdater() *LoadUpdater {
	var lu LoadUpdater
	lu.loadReporter = NewLoadReporter()
	lu.redisWrapper = NewRedisWrapper()

	return &lu
}

func (lu LoadUpdater) CalculateAndAddNewLoad() {
	lastN := 15
	lastMinuteLoad, now := lu.loadReporter.GetAverageLoad(lastN)
	lu.saveDivvyLoad(lastMinuteLoad, now)
}

func (dw LoadUpdater) saveDivvyLoad(load int, timestamp string) {
	_, err := dw.redisWrapper.redisConn.Do("HSET", loadKeys, timestamp, load)
	if err != nil {
		log.Print("Error: ", err)
	}
}

func (dw LoadUpdater) UpdateStationStatuses() {
	ds := dw.divvyApiClient.getCurrentStatuses()
	_, err := dw.redisWrapper.redisConn.Do("ZADD", statusKeys, time.Now().Unix(), ds.ExecutionTime)
	if err != nil {
		log.Print("Error: ", err)
	}
	for _, station := range ds.StationBeanList {
		// Pipeline commands
		dw.redisWrapper.redisConn.Send("HSET", ds.ExecutionTime, station.Id, station.AvailableDocks)
	}
	dw.redisWrapper.redisConn.Flush()
}
