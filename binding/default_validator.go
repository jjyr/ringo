package binding

import (
	"reflect"
	"sync"

	"gopkg.in/go-playground/validator.v8"
)

type defaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

var (
	_      StructValidator = &defaultValidator{}
	config                 = &validator.Config{TagName: "validate"}
)

func (v *defaultValidator) ValidateStruct(obj interface{}) error {
	if kindOfData(obj) == reflect.Struct {
		v.lazyinit()
		if err := v.validate.Struct(obj); err != nil {
			return error(err)
		}
	}
	return nil
}

func (v *defaultValidator) lazyinit() {
	v.once.Do(func() {
		v.validate = validator.New(config)
	})
}

func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()
	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}
