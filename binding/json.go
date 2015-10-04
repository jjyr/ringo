package binding

import (
	"encoding/json"
	"net/http"
)

type jsonBinder struct {
}

func (jsonBinder) Bind(req *http.Request, obj interface{}) error {
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return nil
}
