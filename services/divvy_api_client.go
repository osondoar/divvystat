package services

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/osondoar/divvystat/models"
)

type DivvyApiClient struct {
}

func (dac DivvyApiClient) getCurrentStatuses() models.DivvyStatus {
	var dr models.DivvyStatus

	resp, err := http.Get(divvyApiUrl)
	if err != nil {
		log.Print("Error getting divvy status", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&dr); err != nil {
		log.Println("Could not decode body:", err)
	}

	return dr
}
