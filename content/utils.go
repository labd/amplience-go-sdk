package content

import (
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
)

type ErrorResponse struct {
	Inner      error
	StatusCode int
	Errors     []ErrorObject `json:"errors"`
}

func (e *ErrorResponse) Error() string {
	return e.Errors[0].Message
}

// Unwrap is used to make it work with errors.Is, errors.As.
func (e *ErrorResponse) Unwrap() error {
	return e.Inner
}

type ErrorObject struct {
	Entity       string `json:"entity"`
	Property     string `json:"property"`
	InvalidValue string `json:"invalidValue"`
	Message      string `json:"message"`
}

type Link struct {
	Href string `json:"href"`
}

type PageInformation struct {
	Size          int `json:"size"`
	Number        int `json:"number"`
	TotalElements int `json:"totalElements"`
	TotalPages    int `json:"totalPages"`
}

type GenericListResults struct {
	Embedded map[string]interface{} `json:"_embedded"`
	Links    map[string]Link        `json:"_links"`
	Page     PageInformation        `json:"page"`
}

type ArchiveInput struct {
	Version int `json:"version"`
}

func toTimeHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if t != reflect.TypeOf(time.Time{}) {
			return data, nil
		}

		switch f.Kind() {
		case reflect.String:
			return time.Parse(time.RFC3339, data.(string))
		case reflect.Float64:
			return time.Unix(0, int64(data.(float64))*int64(time.Millisecond)), nil
		case reflect.Int64:
			return time.Unix(0, data.(int64)*int64(time.Millisecond)), nil
		default:
			return data, nil
		}
		// Convert it by parsing
	}
}

func decodeStruct(input interface{}, result interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata:   nil,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(toTimeHookFunc()),
		Result:     result,
		TagName:    "json",
	})
	if err != nil {
		return err
	}

	if err := decoder.Decode(input); err != nil {
		return err
	}
	return err
}
