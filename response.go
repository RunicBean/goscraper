package scraper

import (
	"encoding/json"
	"io"
	"net/http"
)

type Response interface {
	GetStatusCode() int
	GetBody() []byte
	GetBodyMap() (map[string]interface{}, error)
}
type response struct {
	response *http.Response
}

func (r *response) GetStatusCode() int {
	return r.response.StatusCode
}

func (r *response) GetBody() []byte {
	b, _ := io.ReadAll(r.response.Body)
	return b
}

func (r *response) GetBodyMap() (map[string]interface{}, error) {
	b := r.GetBody()
	var m map[string]interface{}
	err := json.Unmarshal(b, &m)
	return m, err
}

func NewResponse(resp *http.Response) Response {
	return &response{response: resp}
}
