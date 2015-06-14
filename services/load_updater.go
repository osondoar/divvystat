package services

import (
	"log"
	"strconv"
	"time"
)

const statusKeys = "status-keys"
const loadKey = "loads"
const expireSeconds = 18000

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
	epoch := GetEpoch(timestamp)
	value := timestamp + "__" + strconv.Itoa(load)
	_, err1 := dw.redisWrapper.redisConn.Do("ZADD", loadKey, epoch, value)
	if err1 != nil {
		log.Print("Error: ", err1)
	}
}

func (dw LoadUpdater) UpdateStationStatuses() {
	unixTime := time.Now().Unix()
	ds := dw.divvyApiClient.getCurrentStatuses()
	_, err := dw.redisWrapper.redisConn.Do("ZADD", statusKeys, unixTime, ds.ExecutionTime)
	if err != nil {
		log.Print("Error: ", err)
	}
	for _, station := range ds.StationBeanList {
		// Pipeline commands
		dw.redisWrapper.redisConn.Do("HSET", ds.ExecutionTime, station.Id, station.AvailableDocks)
	}

	// Delete keys older than 30 mins
	dw.redisWrapper.redisConn.Send("ZREMRANGEBYSCORE", statusKeys, 0, unixTime-expireSeconds)
	dw.redisWrapper.redisConn.Send("EXPIRE", ds.ExecutionTime, expireSeconds)

	dw.redisWrapper.redisConn.Flush()
}
