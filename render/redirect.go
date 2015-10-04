package render

import "net/http"

type Redirect struct {
	Code     int
	Location string
	Request  *http.Request
}

func (r *Redirect) Render(w http.ResponseWriter) error {
	writeContentType(w, jsonContentType)
	http.Redirect(w, r.Request, r.Location, r.Code)
	return nil
}
