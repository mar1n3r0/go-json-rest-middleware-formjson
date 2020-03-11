// Package formjson provides Middleware for converting posted x-www-form-urlencoded data into json
package formjson

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"mime"
	"net/http"
	"strings"

	"github.com/vardius/gorouter/v4"
)

type FormError struct {
	Error   error
	Message string
}

//Provides "x-www-form-urlencoded" to "json" conversion middleware for gorouter
func FormJson() gorouter.MiddlewareFunc {
	m := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

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
						conversionError(w)
						return
					}
				}

				//marshal json
				jsonString, err := json.Marshal(jsonMap)
				if err != nil {
					//error marshalling, skip to handler
					conversionError(w)
					return
				}

				//write new body
				r.Body = ioutil.NopCloser(bytes.NewReader([]byte(string(jsonString))))

				//convert content-type header
				r.Header.Set("Content-Type", "application/json")
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}

	return m
}

func conversionError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	err := json.NewEncoder(w).Encode(map[string]string{"Error": "Error Converting Form Data"})
	if err != nil {
		panic(err)
	}
}
