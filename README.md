# Go HTTP API Contract IO

[![Go](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fproxy.golang.org%2Fgo.leoweyr.com%2Fgo-http-api-contract-io%2Fv4%2F%40latest&query=%24.Version&label=go&logo=go&logoColor=white&color=00ADD8)](https://pkg.go.dev/go.leoweyr.com/go-http-api-contract-io/v4)

Consistent HTTP request and response handling for Go servers and clients, with strict JSON decoding, DTO-based validation, and standardized API errors.

```bash
go get go.leoweyr.com/go-http-api-contract-io/v4
```

The library generates a standardized error response body designed for consistent API communication:

```json
{
  "error": {
    "message": "VALIDATION_FAILED",
    "details": {}
  }
}
```
