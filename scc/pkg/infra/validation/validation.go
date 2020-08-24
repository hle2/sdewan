package validation

import (
	"reflect"
	"fmt"
	"sync"
	"github.com/go-playground/validator/v10"
)

type SdewanValidator struct {
	validate *validator.Validate
}

type safeSdewanValidators struct {
	validates map[string]*SdewanValidator
	mux     sync.Mutex
}

var gvalidates = safeSdewanValidators{validates: make(map[string]*SdewanValidator)}

func GetValidator(name string) *SdewanValidator {
	return gvalidates.getValidate(name)
}

// safeSdewanValidators
func (s *safeSdewanValidators) getValidate(name string) *SdewanValidator {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.validates[name] == nil {
		s.validates[name] = &SdewanValidator{
			validate: validator.New(),
		}
	}

	return s.validates[name]
}

// SdewanValidator
func(v *SdewanValidator) Validate(data interface{}) (bool, string) {
	err := v.validate.Struct(data)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return false, reflect.TypeOf(err).String()
		}

		msg := "Input fields check error: ["
		index := 1
		for _, err := range err.(validator.ValidationErrors) {
			fieldMsg := fmt.Sprintf("%s(%s:%s)", err.Field(), err.Tag(), err.Param())
			if index == 1 {
				msg = msg + fieldMsg
			} else {
				msg = msg + ", " + fieldMsg
			}

			index = index + 1
		}

		msg = msg +"]"

		// from here you can create your own error messages in whatever language you wish
		return false, msg
	}

    return true, ""
}

func(v *SdewanValidator) RegisterValidation(tag string, fn validator.Func) error {
	return v.validate.RegisterValidation(tag, fn)
}

func(v *SdewanValidator) RegisterStructValidation(fn validator.StructLevelFunc, types interface{}) {
	v.validate.RegisterStructValidation(fn, types)
}