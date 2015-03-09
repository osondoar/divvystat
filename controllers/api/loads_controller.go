package api_controllers

import (
	"encoding/json"
	"net/http"

	"github.com/osondoar/divvystat/services"
)

type LoadsController struct {
	ApiController
}

func (controller LoadsController) Index(w http.ResponseWriter, r *http.Request) {

	lr := services.NewLoadReporter()
	loads := lr.GetAverageLoads()
	json, _ := json.Marshal(loads)
	controller.Render(w, r, string(json))
}
