package mph

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"braces.dev/errtrace"
)

type Error struct {
	Title     string
	Detail    error
	ErrorCode int // will be put into Response StatusCode
}

func (e *Error) IsFatal() bool {
	return e != nil && e.ErrorCode >= 500 && e.ErrorCode <= 599
}

type errorJSON struct {
	Title     string `json:"title,omitempty"`
	ErrorCode int    `json:"errorCode,omitempty"`
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Detail
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("%s: %s", e.Title, e.Detail)
}

func (e *Error) MarshalJSON() ([]byte, error) {
	if e == nil || e.Detail == nil {
		return errtrace.Wrap2(json.Marshal(nil))
	}

	if !e.IsFatal() {
		panic(errtrace.Errorf("cannot marshal non-fatal errors to JSON, use ResponseError instead"))
	}

	return errtrace.Wrap2(json.Marshal(errorJSON{
		Title:     e.Title,
		ErrorCode: e.ErrorCode,
	}))
}

func (e *Error) UnmarshalJSON(data []byte) error {
	var ej errorJSON
	if err := json.Unmarshal(data, &ej); err != nil {
		return errtrace.Wrap(err)
	}
	e.Title = ej.Title
	e.ErrorCode = ej.ErrorCode
	return nil
}

var _ json.Marshaler = &Error{}
var _ json.Unmarshaler = &Error{}

func (e *Error) ToResponseError() *ResponseError {
	if e == nil || e.Detail == nil || (e.Title == "" && e.Detail.Error() == "") {
		return nil
	}
	if e.IsFatal() {
		panic(errtrace.Errorf("fatal web.Errors cannot be converted to ResponseError"))
	}
	return &ResponseError{
		Title:  e.Title,
		Detail: e.Detail.Error(),
	}
}

type ResponseError struct {
	Title  string `json:"title,omitempty"  db:"-"`
	Detail string `json:"detail,omitempty" db:"-"`
}

var _ error = &ResponseError{}
var _ sql.Scanner = &ResponseError{}
var _ driver.Valuer = ResponseError{}

func ParseResponseError(err string) *ResponseError {
	if err == "" {
		return nil
	}

	titleIndex := strings.Index(err, ": ")
	if titleIndex == -1 {
		return NewResponseErrorText("Error", err)
	}
	return NewResponseErrorText(err[:titleIndex], err[titleIndex+2:])
}

func NewResponseError(title string, detail error) *ResponseError {
	if detail == nil || detail.Error() == "" {
		return nil
	}
	return NewResponseErrorText(title, detail.Error())
}

func NewResponseErrorText(title, detail string) *ResponseError {
	return &ResponseError{Title: title, Detail: detail}
}

func (r *ResponseError) Error() string {
	if r == nil || r.Detail == "" && r.Title == "" {
		return ""
	}
	return fmt.Sprintf("%s: %s", r.Title, r.Detail)
}

func (r ResponseError) String() string {
	return r.Error()
}

func (r *ResponseError) IsEmpty() bool {
	return r == nil || (r.Title == "" && r.Detail == "")
}

func (r *ResponseError) HasSpecificMessage() bool {
	return r != nil && r.Title != fatalEditErrorTitle && r.Detail != editErrorDetail
}

func (r *ResponseError) ToClientError() *Error {
	if r == nil {
		return nil
	}
	return ClientError(r.Title, errtrace.Errorf("%s", r.Detail))
}

func (r *ResponseError) Scan(value any) error {
	errMsg, ok := value.(string)
	if !ok {
		return errtrace.Errorf("cannot scan %T into ResponseError", value)
	}
	if err := ParseResponseError(errMsg); err != nil {
		*r = *err
	}
	return nil
}

func (r ResponseError) Value() (driver.Value, error) {
	return r.Error(), nil
}

func NewError(title string, err error, errorCode int) *Error {
	if err == nil {
		return nil
	}
	return &Error{
		Title:     title,
		Detail:    err,
		ErrorCode: errorCode,
	}
}

func ForbiddenError(detail error) *Error {
	return NewError("Forbidden", detail, http.StatusForbidden)
}

func ClientError(title string, detail error) *Error {
	var wErr *Error
	if errors.As(detail, &wErr) {
		if wErr.IsFatal() {
			panic(errtrace.Errorf("cannot create client error from fatal error"))
		}
	}
	return NewError(title, detail, http.StatusBadRequest)
}

func InternalError(title string, detail error) *Error {
	return NewError(title, detail, http.StatusInternalServerError)
}
