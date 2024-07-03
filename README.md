# goscraper
## Features
- [x] Simple HTTP request
- [ ] data parser embedded extension
## Simple Tutorial
### Quick call to "python-request-like" library
```go
goscraper.Get("https://www.google.com")
// with custom headers
goscraper.Post(
    "https://www.google.com", 
    goscraper.WithHeaders(
        map[string]string{
            "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.93 Safari/537.36"
        },
    ),
    goscraper.WithJson(`{"name": "John Doe", "age": 30}`),
)
// with authorization
goscraper.Get(
    "https://www.google.com", 
    goscraper.WithBearerToken(
        "token",
    ),
)
```
### Chain request appendix
```go
req, _ := goscraper.NewRequest(testUrl, POST)
response, _ := req.WithData(map[string]string{
    "name": "Jane Doe",
    "age":  "25",
}).WithBearerToken("override").Do()
```