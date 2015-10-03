package render

import (
	"encoding/json"
	"net/http"
)

type JSONData struct {
	Content interface{}
}

var jsonContentType = []string{"application/json; charset=utf-8"}

func (data *JSONData) Render(w http.ResponseWriter) error {
	writeContentType(w, jsonContentType)
	err := json.NewEncoder(w).Encode(data.Content)
	return err
}
