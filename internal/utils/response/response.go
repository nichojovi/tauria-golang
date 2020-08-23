package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// Error
type ErrChain struct {
	Message string
	Cause   error
	Fields  map[string]string
	Type    error
}

func (err ErrChain) Error() string {
	bcoz := ""
	fields := ""
	if err.Cause != nil {
		bcoz = fmt.Sprint(" because {", err.Cause.Error(), "}")
		if len(err.Fields) > 0 {
			fields = fmt.Sprintf(" with Fields {%+v}", err.Fields)
		}
	}
	return fmt.Sprint(err.Message, bcoz, fields)
}

func Type(err error) error {
	switch err.(type) {
	case ErrChain:
		return err.(ErrChain).Type
	}
	return nil
}

func toString(m map[string]string) string {
	v, _ := json.Marshal(m)
	return string(v)
}

func (err ErrChain) SetField(key string, value string) ErrChain {
	if err.Fields == nil {
		err.Fields = map[string]string{}
	}
	err.Fields[key] = value
	return err
}

type InvalidError struct {
	message string
}

func (ie *InvalidError) Error() string {
	return ie.message
}

func NewInvalidError(msg string) *InvalidError {
	return &InvalidError{message: msg}
}

func NewInvalidErrorf(msg string, args ...interface{}) *InvalidError {
	return NewInvalidError(fmt.Sprintf(msg, args...))
}

var (
	ErrBadRequest          = errors.New("Bad request")
	ErrForbidden           = errors.New("Forbidden")
	ErrNotFound            = errors.New("Not found")
	ErrInternalServerError = errors.New("Internal server error")
	ErrInvalidRequest      = errors.New("Invalid request")
)

const (
	STATUSCODE_GENERICSUCCESS = "200000"
	STATUSCODE_BADREQUEST     = "400000"
	STATUS_FORBIDDEN          = "403000"
	STATUSCODE_NOT_FOUND      = "404000"
	STATUSCODE_INTERNAL_ERROR = "500000"
)

func GetErrorCode(err error) string {
	switch err.(type) {
	case ErrChain:
		errType := err.(ErrChain).Type
		if errType != nil {
			err = errType
		}
	}
	switch err {
	case ErrBadRequest:
		return STATUSCODE_BADREQUEST
	case ErrForbidden:
		return STATUS_FORBIDDEN
	case ErrNotFound:
		return STATUSCODE_NOT_FOUND
	case ErrInternalServerError:
		return STATUSCODE_INTERNAL_ERROR
	case ErrInvalidRequest:
		return STATUSCODE_BADREQUEST
	case nil:
		return STATUSCODE_GENERICSUCCESS
	default:
		return STATUSCODE_INTERNAL_ERROR
	}
}

func GetHTTPCode(code string) int {
	s := code[0:3]
	i, _ := strconv.Atoi(s)
	return i
}

// Response
type JSONResponse struct {
	Data        interface{}            `json:"data,omitempty"`
	Message     string                 `json:"message,omitempty"`
	Code        string                 `json:"code"`
	StatusCode  int                    `json:"-"`
	ErrorString string                 `json:"error,omitempty"`
	Error       error                  `json:"-"`
	RealError   string                 `json:"-"`
	Latency     string                 `json:"latency"`
	Log         map[string]interface{} `json:"-"`
}

func NewJSONResponse() *JSONResponse {
	return &JSONResponse{Code: STATUSCODE_GENERICSUCCESS, StatusCode: GetHTTPCode(STATUSCODE_GENERICSUCCESS), Log: map[string]interface{}{}}
}

func (r *JSONResponse) SetData(data interface{}) *JSONResponse {
	r.Data = data
	return r
}

func (r *JSONResponse) SetMessage(msg string) *JSONResponse {
	r.Message = msg
	return r
}

func (r *JSONResponse) SetLatency(latency float64) *JSONResponse {
	r.Latency = fmt.Sprintf("%.2f ms", latency)
	return r
}

func (r *JSONResponse) SetLog(key string, val interface{}) *JSONResponse {
	r.Log[key] = val
	return r
}

func getErrType(err error) error {
	switch err.(type) {
	case ErrChain:
		errType := err.(ErrChain).Type
		if errType != nil {
			err = errType
		}
	}
	return err
}

func (r *JSONResponse) SetError(err error, a ...string) *JSONResponse {
	r.RealError = fmt.Sprintf("%+v", err)
	err = getErrType(err)
	r.Error = err
	r.ErrorString = err.Error()
	r.Code = GetErrorCode(err)
	r.StatusCode = GetHTTPCode(r.Code)
	if r.StatusCode == http.StatusInternalServerError {
		r.ErrorString = "Internal Server error"
	}
	if len(a) > 0 {
		r.ErrorString = a[0]
	}
	return r
}

func (r *JSONResponse) Send(w http.ResponseWriter) {
	b, _ := json.Marshal(r)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.StatusCode)
	w.Write(b)
}
