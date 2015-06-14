package services

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/osondoar/divvystat/models"
)

const divvyApiUrl = "http://www.divvybikes.com/stations/json"
const divvyTimeLayout = "2006-01-02 03:04:05 PM"
const outTimeLayout = "2006-01-02T15:04:05-07:00"

var timeLocation, _ = time.LoadLocation("America/Chicago")

type DivvyApiClient struct {
}

func (dac DivvyApiClient) getCurrentStatuses() models.StationsStatusApi {
	var stationsStatus models.StationsStatusApi

	resp, err := http.Get(divvyApiUrl)
	if err != nil {
		log.Print("Error getting divvy status", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&stationsStatus); err != nil {
		log.Println("Could not decode body:", err)
	}
	stationsStatus.ExecutionTime = dac.parseDivvyTime(stationsStatus.ExecutionTime)
	return stationsStatus
}

func (dac DivvyApiClient) parseDivvyTime(divvyTime string) string {
	t, _ := time.ParseInLocation(divvyTimeLayout, divvyTime, timeLocation)
	return t.Format(outTimeLayout)
}
