package main

import (
	"time"

	"github.com/osondoar/divvystat/services"
)

func main() {
	loadsService := services.NewLoadsService()

	for {
		loadsService.UpdateStationStatuses()
		loadsService.CalculateAndAddNewLoad()
		time.Sleep(60 * time.Second)
	}

}
