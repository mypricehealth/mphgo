package mph

import (
	"encoding/json"
)

// Response contains the standardized API response data used by all My Price Health API's. It is based off of the generalized error handling recommendation found
// in IETF RFC 7807 https://tools.ietf.org/html/rfc7807 and is a simplification of the Spring Boot error response as described at https://www.baeldung.com/rest-api-error-handling-best-practices
/*   An error response might look like this:
{
	"error: {
		"title": "Incorrect username or password.",
		"detail": "Authentication failed due to incorrect username or password.",
	}
	"status": 401,
}

A successful response with a single result might look like this:
{
	"result": {
		"procedureCode": "ABC",
		"billedAverage": 15.23
	},
	"status": 200,
}
*/
type Response[Result any] struct {
	Error      *ResponseError `json:"error,omitempty"`  // supplied when entire response is an error
	Result     Result         `json:"result,omitempty"` // supplied on success. Will be a single object.
	StatusCode int            `json:"status"`           // supplied on success and error
}

type responseJSON[Result any] struct {
	Message    string         `json:"message,omitempty"`
	Code       int            `json:"code,omitempty"`
	Error      *ResponseError `json:"error,omitempty"`
	Result     Result         `json:"result,omitempty"`
	StatusCode int            `json:"status"`
}

func (r *Response[Result]) UnmarshalJSON(data []byte) error {
	var rj responseJSON[Result]
	if err := json.Unmarshal(data, &rj); err != nil {
		return err
	}
	r.Error = rj.Error
	r.Result = rj.Result
	r.StatusCode = rj.StatusCode
	if rj.Code != 0 {
		r.StatusCode = rj.Code
		r.Error = &ResponseError{Title: rj.Message}
	}
	return nil
}

// Responses contains the standardized API response data used by all My Price Health API's. It is based off of the generalized error handling recommendation found
// in IETF RFC 7807 https://tools.ietf.org/html/rfc7807 and is a simplification of the Spring Boot error response as described at https://www.baeldung.com/rest-api-error-handling-best-practices
/*   A response with one success and one failure might look like this:
{
	"results": [
		{
			"procedureCode": "ABC",
			"billedAverage": 15.23
		},
		{
			"error": {
				"title": "invalid procedure code",
				"detail": "unable to find procedure code `DEF` in the list of valid procedure codes"
			}
		}
	],
	"status": 200,
	"successCount": 1,
	"errorCount": 1,
}

A successful response with multiple results might look like this (note no embedded error field):
{
	"results": [
		{
			"procedureCode": "ABC",
			"billedAverage": 15.23
		},
		{
			"procedureCode": "DEF",
			"billedAverage": 12.22
		}
	],
	"status": 200,
	"successCount": 2,
	"errorCount": 0,
}
*/
type Responses[Result any] struct {
	Error        *ResponseError `json:"error,omitempty"`   // supplied when entire response is an error
	Results      []Result       `json:"results,omitempty"` // A slice of results that will either be a successful result or an error.
	SuccessCount int            `json:"successCount"`      // count of successful results when WriteResults is called
	ErrorCount   int            `json:"errorCount"`        // count of errored results when WriteResults is called
	StatusCode   int            `json:"status"`            // supplied on success and error
}

type responsesJSON[Result any] struct {
	Message      string         `json:"message,omitempty"`
	Code         int            `json:"code,omitempty"`
	Error        *ResponseError `json:"error,omitempty"`
	Results      []Result       `json:"results,omitempty"`
	SuccessCount int            `json:"successCount"`
	ErrorCount   int            `json:"errorCount"`
	StatusCode   int            `json:"status"`
}

func (r *Responses[Result]) UnmarshalJSON(data []byte) error {
	var rj responsesJSON[Result]
	if err := json.Unmarshal(data, &rj); err != nil {
		return err
	}
	r.Error = rj.Error
	r.Results = rj.Results
	r.SuccessCount = rj.SuccessCount
	r.ErrorCount = rj.ErrorCount
	r.StatusCode = rj.StatusCode
	if rj.Code != 0 {
		r.StatusCode = rj.Code
		r.Error = &ResponseError{Title: rj.Message}
	}
	return nil
}

// ResponseError supplies detailed error information when an entire request or an item in a response fails
type ResponseError struct {
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
}

func (e *ResponseError) Error() string {
	return e.Title + ": " + e.Detail
}
