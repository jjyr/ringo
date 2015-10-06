package render

import "net/http"

type Redirect struct {
	Code     int
	Location string
	Request  *http.Request
}

func (r *Redirect) Render(w http.ResponseWriter) error {
	http.Redirect(w, r.Request, r.Location, r.Code)
	return nil
}
