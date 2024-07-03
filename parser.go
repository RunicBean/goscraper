package goscraper

import (
	"net/http"
)

type ParseFunc func(*http.Response) (interface{}, error)
