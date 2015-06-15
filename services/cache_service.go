package services

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/osondoar/divvystat/models"
)

const statusKeys = "status-keys"
const loadKey = "loads"
const expireSeconds = 18000

type CacheService struct {
	host  string
	redis redis.Conn
}

func NewCacheService() *CacheService {
	var cs CacheService
	var err error

	hostEnv := os.Getenv("REDIS_PORT_6379_TCP_ADDR")
	if hostEnv != "" {
		cs.host = hostEnv
	} else {
		cs.host = "localhost"
	}

	cs.redis, err = redis.Dial("tcp", fmt.Sprintf("%s:6379", cs.host))
	if err != nil {
		log.Fatal("Can't connect to redis: ", err)
	}

	return &cs
}

func (cs CacheService) getStatuses(lastN int) []models.StationsStatus {
	var lastStatuses []models.StationsStatus

	lastStatusKeys, err := redis.Strings(cs.redis.Do("ZREVRANGE", statusKeys, 0, lastN))
	if err != nil {
		log.Print("Error: ", err)
	}

	for _, statusKey := range lastStatusKeys {
		stationsStatus := models.NewStationsStatus(statusKey)

		statusTuples, err1 := redis.Ints(cs.redis.Do("HGETALL", statusKey))
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

func (cs CacheService) addStatuses(statuses models.StationsStatusApi) {
	unixTime := time.Now().Unix()

	_, err := cs.redis.Do("ZADD", statusKeys, unixTime, statuses.ExecutionTime)
	if err != nil {
		log.Print("Error: ", err)
	}
	for _, station := range statuses.StationBeanList {
		cs.redis.Do("HSET", statuses.ExecutionTime, station.Id, station.AvailableDocks)
	}

	cs.redis.Send("ZREMRANGEBYSCORE", statusKeys, 0, unixTime-expireSeconds)
	cs.redis.Send("EXPIRE", statuses.ExecutionTime, expireSeconds)

	cs.redis.Flush()
}

func (cs CacheService) getEncodedLoads(from int64, to int64) []string {
	encodedLoads, err1 := redis.Strings(cs.redis.Do("ZREVRANGEBYSCORE", loadKey, to, from))

	if err1 != nil {
		log.Print("Error: ", err1)
	}

	return encodedLoads
}

func (cs CacheService) addLoad(epoch int64, encodedValue string) {
	_, err1 := cs.redis.Do("ZADD", loadKey, epoch, encodedValue)
	if err1 != nil {
		log.Print("Error: ", err1)
	}
}
