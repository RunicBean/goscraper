package scraper

import (
	"net/http"
)

type ParseFunc func(*http.Response) (any, error)
