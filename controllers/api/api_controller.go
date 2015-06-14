package api_controllers

import (
	"encoding/json"
	"net/http"
)

type ApiController struct {
}

func (controller ApiController) Render(w http.ResponseWriter, r *http.Request, body interface{}) {
	json, _ := json.Marshal(body)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(string(json)))
}
