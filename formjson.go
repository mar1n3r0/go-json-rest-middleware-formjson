// Package formjson provides Middleware for converting posted x-www-form-urlencoded data into json
package formjson

import (
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

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				err := json.NewEncoder(w).Encode(jsonMap)
				if err != nil {
					panic(err)
				}

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
