package formjson

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

func TestFormJson(t *testing.T) {
	m := FormJson()
	h := m(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))

	data := url.Values{}
	data.Set("name", "foo")
	data.Add("surname", "bar")

	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/", bytes.NewBufferString(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	if err != nil {
		t.Fatal(err)
	}

	h.ServeHTTP(w, req)

	if w.Header().Get("Content-Type") != "application/json" {
		t.Error("FormJson returned unexpected headers: ", w.Header())
	}

	if w.Code != http.StatusOK {
		t.Error("Test http form request returned unexpected status code: ", w.Result().StatusCode)
	}

	cmp := bytes.Compare(w.Body.Bytes(), append([]byte(`{"name":"foo","surname":"bar"}`), 10))
	if cmp != 0 {
		t.Errorf("FormJson returned unexpected body: %s | %d", w.Body.String(), cmp)
	}
}

func TestFormJsonError(t *testing.T) {
	m := FormJson()
	h := m(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))

	data := "mal&formed"

	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/", bytes.NewBufferString(data))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data)))
	if err != nil {
		t.Fatal(err)
	}

	h.ServeHTTP(w, req)

	if w.Header().Get("Content-Type") != "application/json" {
		t.Error("FormJson returned unexpected headers: ", w.Header())
	}

	if w.Code != http.StatusInternalServerError {
		t.Error("Test http form request returned unexpected status code: ", w.Result().StatusCode)
	}

	cmp := bytes.Compare(w.Body.Bytes(), append([]byte(`{"Error":"Error Converting Form Data"}`), 10))
	if cmp != 0 {
		t.Errorf("FormJson returned unexpected body: %s | %d", w.Body.String(), cmp)
	}
}
