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

	loadReporter := services.NewLoadReporter()

	Last1, _ := loadReporter.GetAverageLoad(1)
	Last5, _ := loadReporter.GetAverageLoad(5)
	Last15, _ := loadReporter.GetAverageLoad(10)
	Last60, _ := loadReporter.GetAverageLoad(60)
	report := &models.Report{Last1, Last5, Last15, Last60}
	controller.RenderTemplate(w, r, "index.html", report)

}
