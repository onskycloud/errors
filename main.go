// Package errors provides a way to return detailed information
// for an RPC request error. The error is normally JSON encoded.
package errors

import (
	"encoding/json"
	"fmt"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	// RedisEmpty redis empty value response
	RedisEmpty = "redis: nil"
	// NotSupport means language not support yet
	NotSupport = "language:notSupport"
	// NotExisted means message type not existed
	NotExisted = "messageType:notExisted"
)

// Error implements the error interface.
type Error struct {
	Id     string `json:"id"`
	Code   int32  `json:"code"`
	Detail string `json:"detail"`
	Status string `json:"status"`
}

// ErrorDict represents error list
type ErrorDict struct {
	ErrorList []ErrorMessage `yaml:"error_list,omitempty" json:"errorList,omitempty"`
}

// ErrorMessage holds message type for many languages
type ErrorMessage struct {
	Type              string    `yaml:"type,omitempty" json:"type,omitempty"`
	TranslatedMessage []Message `yaml:"translated_message,omitempty" json:"translated_message,omitempty"`
}

// Message represents message for a language
type Message struct {
	Text     string `yaml:"text,omitempty" json:"text,omitempty"`
	Language string `yaml:"language,omitempty" json:"language,omitempty"`
}

func (e *Error) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

// New generates a custom error.
func New(id, detail string, code int32) error {
	return &Error{
		Id:     id,
		Code:   code,
		Detail: detail,
		Status: http.StatusText(int(code)),
	}
}

// Parse tries to parse a JSON string into an error. If that
// fails, it will set the given string as the error detail.
func Parse(err string) *Error {
	e := new(Error)
	errr := json.Unmarshal([]byte(err), e)
	if errr != nil {
		e.Detail = err
	}
	return e
}

// BadRequest generates a 400 error.
func BadRequest(id, format string, a ...interface{}) error {
	return &Error{
		Id:     id,
		Code:   400,
		Detail: fmt.Sprintf(format, a...),
		Status: http.StatusText(400),
	}
}

// Unauthorized generates a 401 error.
func Unauthorized(id, format string, a ...interface{}) error {
	return &Error{
		Id:     id,
		Code:   401,
		Detail: fmt.Sprintf(format, a...),
		Status: http.StatusText(401),
	}
}

// Forbidden generates a 403 error.
func Forbidden(id, format string, a ...interface{}) error {
	return &Error{
		Id:     id,
		Code:   403,
		Detail: fmt.Sprintf(format, a...),
		Status: http.StatusText(403),
	}
}

// NotFound generates a 404 error.
func NotFound(id, format string, a ...interface{}) error {
	return &Error{
		Id:     id,
		Code:   404,
		Detail: fmt.Sprintf(format, a...),
		Status: http.StatusText(404),
	}
}

// MethodNotAllowed generates a 405 error.
func MethodNotAllowed(id, format string, a ...interface{}) error {
	return &Error{
		Id:     id,
		Code:   405,
		Detail: fmt.Sprintf(format, a...),
		Status: http.StatusText(405),
	}
}

// Timeout generates a 408 error.
func Timeout(id, format string, a ...interface{}) error {
	return &Error{
		Id:     id,
		Code:   408,
		Detail: fmt.Sprintf(format, a...),
		Status: http.StatusText(408),
	}
}

// Conflict generates a 409 error.
func Conflict(id, format string, a ...interface{}) error {
	return &Error{
		Id:     id,
		Code:   409,
		Detail: fmt.Sprintf(format, a...),
		Status: http.StatusText(409),
	}
}

// InternalServerError generates a 500 error.
func InternalServerError(id, format string, a ...interface{}) error {
	return &Error{
		Id:     id,
		Code:   500,
		Detail: fmt.Sprintf(format, a...),
		Status: http.StatusText(500),
	}
}

// ErrorMessageTranslater converts messageType to message for input language
func ErrorMessageTranslater(path string, messageType string, language string) (string, error) {
	var transMsg string
	errorDict, err := LoadErrorList(path)
	if err != nil {
		return "", err
	}
	if errorDict.ErrorList != nil && len(errorDict.ErrorList) > 0 {
		for _, v := range errorDict.ErrorList {
			if v.Type == messageType {
				for _, t := range v.TranslatedMessage {
					if t.Language == language {
						transMsg = t.Text
						return transMsg, nil
					}
				}
				return NotSupport, nil
			}
		}
		return NotExisted, nil
	}

	return "", nil
}

// LoadErrorList loads error list from file
func LoadErrorList(path string) (*ErrorDict, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading error list file, %s", err)
	}
	var cfg = new(ErrorDict)
	if err := yaml.Unmarshal(bytes, cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)
	}
	return cfg, nil
}
