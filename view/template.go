package view

import (
	"html/template"
	"net/http"
	"path"
)

func fetchTemplates(w http.ResponseWriter, name string, data any) {
	var (
		tmpl *template.Template
		err  error
	)
	tmpl, err = template.ParseFiles(path.Join("templates", name+".html"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
