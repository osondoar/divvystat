package controllers

import (
	"net/http"
	"path"
	"text/template"
)

type AppController struct {
}

func (controller AppController) Render(w http.ResponseWriter, r *http.Request, body string) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(body))
}

func (controller AppController) RenderTemplate(w http.ResponseWriter, r *http.Request, templateName string, params interface{}) {
	w.Header().Set("Content-Type", "text/html")

	fp := path.Join("templates", templateName)
	template, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := template.Execute(w, params); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
