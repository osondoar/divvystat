package controllers

import (
	"net/http"

	"github.com/osondoar/divvystat/models"
	"github.com/osondoar/divvystat/services"
)

type MainController struct {
	AppController
}

func (controller MainController) Index(w http.ResponseWriter, r *http.Request) {

	loadsService := services.NewLoadsService()

	Last1, _ := loadsService.CurrentAveragedLoad(1)
	Last5, _ := loadsService.CurrentAveragedLoad(5)
	Last15, _ := loadsService.CurrentAveragedLoad(10)
	Last60, _ := loadsService.CurrentAveragedLoad(60)
	report := &models.Report{Last1, Last5, Last15, Last60}
	controller.RenderTemplate(w, r, "index.html", report)

}
