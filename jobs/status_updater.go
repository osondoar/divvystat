package main

import (
	"time"

	"github.com/osondoar/divvystat/services"
)

func main() {
	load_updater := services.NewLoadUpdater()

	for {
		load_updater.UpdateStationStatuses()
		load_updater.CalculateAndAddNewLoad()
		time.Sleep(60 * time.Second)
	}

}
