package goscraper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Response interface {
	GetStatusCode() int
	GetBody() ([]byte, error)
	GetBodyMap() (map[string]interface{}, error)
	Unmarshal(v interface{}) error
}
type response struct {
	response *http.Response
}

func (r *response) GetStatusCode() int {
	return r.response.StatusCode
}

func (r *response) GetBody() ([]byte, error) {
	defer r.response.Body.Close()
	b, err := io.ReadAll(r.response.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (r *response) GetBodyMap() (map[string]interface{}, error) {
	b, err := r.GetBody()
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	err = json.Unmarshal(b, &m)
	return m, err
}

func (r *response) Unmarshal(v interface{}) error {
	b, err := r.GetBody()
	if err != nil {
		return fmt.Errorf("unmarshal: %s", err)
	}
	return json.Unmarshal(b, v)
}

func NewResponse(resp *http.Response) Response {
	return &response{response: resp}
}
