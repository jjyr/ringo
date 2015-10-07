package render

import (
	"html/template"
	"net/http"
)

type HTMLTemplateData struct {
	Template *template.Template
	Data     interface{}
}

var htmlContentType = []string{"text/html; charset=utf-8"}

func (htmlTemplateData *HTMLTemplateData) Render(w http.ResponseWriter) error {
	writeContentType(w, htmlContentType)
	err := htmlTemplateData.Template.Execute(w, htmlTemplateData.Data)
	return err
}
