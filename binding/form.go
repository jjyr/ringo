package binding

import "net/http"

type formBinder struct{}
type formPostBinder struct{}
type formMultipartBinder struct{}

func (formBinder) Bind(req *http.Request, obj interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	req.ParseMultipartForm(32 << 10) // 32 MB
	if err := mapForm(obj, req.Form); err != nil {
		return err
	}
	return nil
}

func (formPostBinder) Bind(req *http.Request, obj interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	if err := mapForm(obj, req.PostForm); err != nil {
		return err
	}
	return nil
}

func (formMultipartBinder) Bind(req *http.Request, obj interface{}) error {
	if err := req.ParseMultipartForm(32 << 10); err != nil {
		return err
	}
	if err := mapForm(obj, req.MultipartForm.Value); err != nil {
		return err
	}
	return nil
}
