package render

import (
	"fmt"
	"net/http"
)

type TextData struct {
	Format   string
	Contents []interface{}
}

var textContentType = []string{"text/plain; charset=utf-8"}

func (textData *TextData) Render(w http.ResponseWriter) error {
	writeContentType(w, textContentType)
	byteContent := []byte(fmt.Sprintf(textData.Format, textData.Contents...))
	_, err := w.Write(byteContent)
	return err
}
