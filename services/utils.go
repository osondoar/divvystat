package services

import (
	"log"
	"time"
)

func GetEpoch(timestamp string) int64 {
	time, err := time.ParseInLocation(outTimeLayout, timestamp, timeLocation)
	if err != nil {
		log.Print("Error parsing timestamp: ", timestamp, err)
	}
	return time.Unix()

}
