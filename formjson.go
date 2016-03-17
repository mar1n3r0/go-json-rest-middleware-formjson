// Package formjson provides Middleware for converting posted x-www-form-urlencoded data into json
package formjson

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"mime"
	"net/http"
	"strings"

	"github.com/ant0ine/go-json-rest/rest"
)

//Middleware for converting posted form data into json
type Middleware struct{}

func (mw *Middleware) MiddlewareFunc(handler rest.HandlerFunc) rest.HandlerFunc {

	return func(w rest.ResponseWriter, r *rest.Request) {

		mediatype, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if mediatype == "application/x-www-form-urlencoded" {
			// get body
			buf, _ := ioutil.ReadAll(r.Body)

			// map body form data
			jsonMap := map[string]string{}
			sections := strings.Split(string(buf), "&")
			for _, sectionValue := range sections {
				sectionParts := strings.Split(sectionValue, "=")
				if len(sectionParts) == 2 {
					jsonMap[sectionParts[0]] = sectionParts[1]
				} else {
					//error converting, skip to handler
					mw.conversionError(w)
					return
				}
			}

			//marshal json
			jsonString, err := json.Marshal(jsonMap)
			if err != nil {
				//error marshalling, skip to handler
				mw.conversionError(w)
				return
			}

			//write new body
			r.Body = ioutil.NopCloser(bytes.NewReader([]byte(string(jsonString))))

			//convert content-type header
			r.Header.Set("Content-Type", "application/json")
		}
		// call the wrapped handle
		handler(w, r)
	}
}

func (mw *Middleware) conversionError(w rest.ResponseWriter) {
	rest.Error(w, "Error Converting Form Data", http.StatusInternalServerError)
}
