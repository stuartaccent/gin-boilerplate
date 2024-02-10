package webx

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
)

// GetFieldErrors parses validation errors and returns a map keyed by
// the field name with a friendly error message.
//
// Usage:
//
//	if err := c.ShouldBind(&mystruct); err != nil {
//		var ve validator.ValidationErrors
//		if errors.As(err, &ve) {
//			log.Print(webx.GetFieldErrors(ve))
//		}
//	}
func GetFieldErrors(ve validator.ValidationErrors) map[string]string {
	out := make(map[string]string)
	for _, fe := range ve {
		var msg string
		switch fe.Tag() {
		case "email":
			msg = "Invalid email address"
		case "max":
			if fe.Kind() == reflect.String {
				msg = fmt.Sprintf("Must be at most %v characters long", fe.Param())
			} else {
				msg = fmt.Sprintf("Must be at most %v", fe.Param())
			}
		case "min":
			if fe.Kind() == reflect.String {
				msg = fmt.Sprintf("Must be at least %v characters long", fe.Param())
			} else {
				msg = fmt.Sprintf("Must be at least %v", fe.Param())
			}
		case "required":
			msg = "This field is required"
		default:
			msg = fe.Tag()
		}
		out[fe.Field()] = msg
	}
	return out
}
