#go-json-rest-middleware-formjson
x-www-form-urlencoded to json conversion middleware for go-json-rest

*Currently offers basic functionality, API subject to change*

## Explanation

This package provides a [Go-Json-Rest](https://ant0ine.github.io/go-json-rest/) middleware useful converting request data with the content type "application/x-www-form-urlencoded" and converting it to "application/json"

Used BEFORE ContentTypeCheckerMiddleware in cases where x-www-form-urlencoded data must be handled

If "Content-Type" Header set to "x-www-form-urlencoded":

* Changes header "Content-Type" from "application/x-www-form-urlencoded" to "application/json"
* Converts body from `var1=val1&var2=val2` to `{"var1":"val1","var2":"val2"}`

## Installation

    go get github.com/boonep/go-json-rest-middleware-formjson

## Usage

Example 1

Check for Content-Type application/x-www-form-urlencoded on all routes, and convert to application/json

	formJsonMiddleware := &formJson.MiddleWare{}

	api.Use([]rest.Middleware{
		formJsonMiddleware,
		&rest.ContentTypeCheckerMiddleware{},
	}...)

Example 2

Only check specific route and request Type.  This will be the most likely use case.  You should use json content wherever possible, but if you MUST interact with form data on a specific endpoint, you can handle it this way.


	formJsonMiddleware := &formJson.MiddleWare{}

	api.Use(&rest.IfMiddleware{
		Condition: func(request *rest.Request) bool {
			return request.URL.Path == "/form-data" && request.Method == "POST"
		},
		IfTrue: formJsonMiddleware,
	})

	api.Use([]rest.Middleware{
		&rest.ContentTypeCheckerMiddleware{},
	}...)

Above will only convert data POSTed to the "/form-data" path

## Notes

This middleware performs basic functionality for our specific use case (interact with a 3rd party SAAS solution that could not provide JSON content).  Feel free to improve and submit pull requests.  Thanks!




