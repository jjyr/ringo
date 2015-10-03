package render

import "net/http"

type Renderable interface {
	Render(http.ResponseWriter) error
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}
