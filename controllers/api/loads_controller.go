package api_controllers

import (
	"math"
	"net/http"
	"time"

	"github.com/osondoar/divvystat/services"
)

type LoadsController struct {
	ApiController
}

func (controller LoadsController) Index(w http.ResponseWriter, r *http.Request) {
	fromParam := r.URL.Query().Get("from")
	toParam := r.URL.Query().Get("to")
	var from, to int64
	if fromParam != "" {
		from = services.GetEpoch(fromParam)
	} else {
		from = time.Now().AddDate(0, 0, -7).Unix()
	}

	if toParam != "" {
		to = services.GetEpoch(toParam)
	} else {
		to = math.MaxInt64
	}

	lr := services.NewLoadReporter()
	loads := lr.GetAverageLoads(from, to)

	controller.Render(w, r, loads)
}
