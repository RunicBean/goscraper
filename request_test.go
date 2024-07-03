package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
)

func TestCommonCall(t *testing.T) {
	var testUrl = "http://headers.jsontest.com/"
	req, err := NewRequest(testUrl, GET)
	if err != nil {
		t.Fatalf("Error creating request: %s", err)
	}
	resp, err := req.Do()
	if err != nil {
		t.Fatalf("Error getting response: %s", err)
	}
	if resp.GetStatusCode() != 200 {
		t.Fatalf("Expected status code 200, got %d", resp.GetStatusCode())
	}
	var bodyMap map[string]string
	if err := json.Unmarshal(resp.GetBody(), &bodyMap); err != nil {
		t.Fatalf("Error unmarshalling response body: %s", err)
	}
	if bodyMap["Host"] != "headers.jsontest.com" {
		t.Fatalf("Expected Host to be headers.jsontest.com, got %s", bodyMap["Host"])
	}
}

func TestOptions(t *testing.T) {
	testUrl := "https://httpbin.org/bearer"
	req, err := NewRequest(testUrl, GET,
		WithHeaders(
			map[string]string{
				"Accept": "application/json",
				"Host":   "reqbin.com",
			}),
		WithBearerToken("{token}"),
	)
	if err != nil {
		t.Fatalf("Error creating request: %s", err)
	}
	req.WithBearerToken("override")
	resp, err := req.Do()
	if err != nil {
		t.Fatalf("Error getting response: %s", err)
	}
	if resp.GetStatusCode() != 200 {
		t.Fatalf("Expected status code 200, got %d", resp.GetStatusCode())
	}
	bm, err := resp.GetBodyMap()
	if err != nil {
		t.Fatalf("Error getting body map: %s", err)
	}
	if bm == nil || bm["authenticated"] == false {
		t.Fatalf("Expected authenticated to be true, got %v", bm["success"])
	}
}

func TestForbidden(t *testing.T) {
	testUrl := "https://httpbin.org/status/403"
	req, err := NewRequest(testUrl, GET)
	if err != nil {
		t.Fatalf("Error creating request: %s", err)
	}
	resp, err := req.Do()
	if err != nil {
		t.Fatalf("Error getting response: %s", err)
	}
	if resp.GetStatusCode() != 403 {
		t.Fatalf("Expected status code 403, got %d", resp.GetStatusCode())
	}
}

func TestWithData(t *testing.T) {
	testUrl := "https://httpbin.org/post"
	req, err := NewRequest(testUrl, POST,
		WithData(map[string]string{
			"name": "John Doe",
			"age":  "30",
		}),
	)
	if err != nil {
		t.Fatalf("Error creating request: %s", err)
	}
	req.WithData(map[string]string{
		"name": "Jane Doe",
		"age":  "25",
	})
	resp, err := req.Do()
	if err != nil {
		t.Fatalf("Error getting response: %s", err)
	}
	d, err := resp.GetBodyMap()
	if err != nil {
		t.Fatalf("Error getting body map: %s", err)
	}
	if d["form"].(map[string]interface{})["age"] != "25" {
		t.Fatalf("Expected age to be 25, got form %s", d["form"])
	}
	fmt.Println(d)
}

func TestWithJson(t *testing.T) {
	testUrl := "https://httpbin.org/post"
	req, err := NewRequest(testUrl, POST,
		WithJson(`{"name": "John Doe", "age": 30}`),
	)
	if err != nil {
		t.Fatalf("Error creating request: %s", err)
	}
	req.WithJson(`{"name": "Jane Doe", "age": 25}`)
	resp, err := req.Do()
	if err != nil {
		t.Fatalf("Error getting response: %s", err)
	}
	d, err := resp.GetBodyMap()
	if err != nil {
		t.Fatalf("Error getting body map: %s", err)
	}
	fmt.Println(d["json"].(map[string]interface{})["age"])
	if d["json"].(map[string]interface{})["age"] != float64(25) {
		t.Fatalf("Expected age to be 25, got json %s", d["json"])
	}
}

func TestTimeout(t *testing.T) {
	testUrl := "https://httpbin.org/delay/5"
	ctx, cancel := context.WithTimeout(context.Background(), 2)
	defer cancel()
	req, err := NewRequest(testUrl, GET,
		WithContext(ctx),
	)
	if err != nil {
		t.Fatalf("Error creating request: %s", err)
	}
	req.WithContext(ctx)
	_, err = req.Do()
	if err == nil {
		t.Fatalf("Expected timeout error, but got nil")
	}
}

func TestQuickCall(t *testing.T) {
	testUrl := "https://httpbin.org/get"
	resp, err := Get(testUrl)
	if err != nil {
		t.Fatalf("Error test quick call get: %s", err)
	}
	if resp.GetStatusCode() != 200 {
		t.Fatalf("Expected status code 200, but got %d", resp.GetStatusCode())
	}

	testUrl = "https://httpbin.org/post"
	resp, err = Post(testUrl, WithData(map[string]string{"data": "test"}))
	if err != nil {
		t.Fatalf("Error test quick call post: %s", err)
	}
	if resp.GetStatusCode() != 200 {
		t.Fatalf("Expected status code 200, but got %d", resp.GetStatusCode())
	}
}
