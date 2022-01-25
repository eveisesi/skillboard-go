package errors

import "github.com/gobuffalo/buffalo"

func NewBuffaloHTTPError(code int, err error) error {
	return buffalo.HTTPError{
		Status: code,
		Cause:  err,
	}
}
