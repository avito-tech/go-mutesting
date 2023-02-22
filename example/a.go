package example

import (
	"net/http"
)

type A struct{}

func fooA() (a A, b http.Header) { //nolint:unused
	_, _, _ = a, b, http.Header{}

	return a, b
}
