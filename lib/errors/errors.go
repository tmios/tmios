package errors

import (
	"encoding/json"
	"fmt"
	"log"
)

type Error struct {
	Code     int
	Status   int
	Detail   string
	Internal string      `json:"Internal,omitempty"`
	Content  interface{} `json:"Content,omitempty"`
}

func (e Error) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

func (e Error) SetDetail(fmtStr string, args ...interface{}) Error {
	newErr := e
	newErr.Detail += fmt.Sprintf(fmtStr, args...)

	return newErr
}

func Parse(err string) Error {
	var e Error
	errr := json.Unmarshal([]byte(err), &e)
	if errr != nil {
		e.Detail = err
	}
	return e
}

var Errors = map[int]Error{}

func addError(err Error) Error {
	e, ok := Errors[err.Code]
	if ok {
		log.Fatalf("duplate error code: %v, %v", e, err)
	}

	Errors[err.Code] = err
	return err
}

func BadRequest(code int, detail string) Error {
	return addError(Error{
		Code:   code,
		Status: 400,
		Detail: detail,
	})
}

func Conflict(code int, detail string) Error {
	return addError(Error{
		Code:   code,
		Status: 409,
		Detail: detail,
	})
}

func Unauthorized(code int, detail string) Error {
	return addError(Error{
		Code:   code,
		Status: 401,
		Detail: detail,
	})
}

func Forbidden(code int, detail string) Error {
	return addError(Error{
		Code:   code,
		Status: 403,
		Detail: detail,
	})
}

func NotFound(code int, detail string) Error {
	return addError(Error{
		Code:   code,
		Status: 404,
		Detail: detail,
	})
}

func Internal(detail string, err error) Error {
	e, ok := err.(Error)
	if ok {
		return e
	}

	internal := ""
	if err != nil {
		internal = err.Error()
	}

	return Error{
		Code:     500,
		Status:   500,
		Detail:   detail,
		Internal: internal,
	}
}

func InternalNew(detail, internal string) Error {
	return Error{
		Code:     500,
		Status:   500,
		Detail:   detail,
		Internal: internal,
	}
}
