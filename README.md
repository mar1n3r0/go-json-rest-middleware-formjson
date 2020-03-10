#gorouter-middleware-formjson
Provides "x-www-form-urlencoded" to "json" conversion middleware for go-json-rest

## Explanation

This package provides a [Gorouter](https://github.com/vardius/gorouter) middleware useful for converting request data with the content type "application/x-www-form-urlencoded" to "application/json"

If "Content-Type" Header set to "x-www-form-urlencoded":

* Changes header "Content-Type" from "application/x-www-form-urlencoded" to "application/json"
* Converts body from `var1=val1&var2=val2` to `{"var1":"val1","var2":"val2"}`

## Installation

    go get github.com/mar1n3r0/gorouter-middleware-formjson

## Usage

```
// NewRouter provides new router
func NewRouter(logger *log.Logger, server *server.Server, mysqlConnection *sql.DB, grpcConnectionMap map[string]*grpc.ClientConn) gorouter.Router {
	// Global middleware
	router := gorouter.New(
		http_form_middleware.FormJson()
	)
```
