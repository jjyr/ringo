//Package binding idea & some code from gin binding https://github.com/gin-gonic/gin/tree/master/binding , thanks gin
package binding

import (
	"net/http"
)

type Binder interface {
	Bind(*http.Request, interface{}) error
}

const (
	MIMEJSON              = "application/json"
	MIMEHTML              = "text/html"
	MIMEXML               = "application/xml"
	MIMEXML2              = "text/xml"
	MIMEPlain             = "text/plain"
	MIMEPOSTForm          = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm = "multipart/form-data"
)

type StructValidator interface {
	// ValidateStruct can receive any kind of type and it should never panic, even if the configuration is not right.
	// If the received type is not a struct, any validation should be skipped and nil must be returned.
	// If the received type is a struct or pointer to a struct, the validation should be performed.
	// If the struct is not valid or the validation itself fails, a descriptive error should be returned.
	// Otherwise nil must be returned.
	ValidateStruct(interface{}) error
}

var Validator StructValidator = &defaultValidator{}

var (
	JsonBinding          = jsonBinder{}
	FormBinding          = formBinder{}
	FormPostBinding      = formPostBinder{}
	FormMultipartBinding = formMultipartBinder{}
)

func Default(method, contentType string) Binder {
	if method == "GET" {
		return FormBinding
	} else {
		switch contentType {
		case MIMEJSON:
			return JsonBinding
		default: //case MIMEPOSTForm, MIMEMultipartPOSTForm:
			return FormBinding
		}
	}
}

func validate(obj interface{}) error {
	if Validator == nil {
		return nil
	}
	return Validator.ValidateStruct(obj)
}

func BindWith(req *http.Request, obj interface{}, b Binder) error {
	err := b.Bind(req, obj)

	if err == nil {
		err = validate(obj)
	}
	return err
}
