package formjson

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

func TestFormJson(t *testing.T) {
	m := FormJson()
	h := m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

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

	body, e := ioutil.ReadAll(req.Body)
	if e != nil {
		panic(e)
	}

	if req.Header.Get("Content-type") != "application/json" {
		t.Error("FormJson returned unexpected headers: ", req.Header)
	}

	if w.Code != http.StatusOK {
		t.Error("Test http form request returned unexpected status code: ", w.Code)
	}

	cmp := bytes.Compare(body, append([]byte(`{"name":"foo","surname":"bar"}`)))
	if cmp != 0 {
		t.Errorf("FormJson returned unexpected body: %s | %d", body, cmp)
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

	if w.Header().Get("Content-type") != "application/json" {
		t.Error("FormJson returned unexpected headers: ", req.Header)
	}

	if w.Code != http.StatusInternalServerError {
		t.Error("Test http form request returned unexpected status code: ", w.Code)
	}

	cmp := bytes.Compare(w.Body.Bytes(), append([]byte(`{"code":"internal","message":"Error converting form data"}`), 10))
	if cmp != 0 {
		t.Errorf("FormJson returned unexpected body: %s | %d", w.Body.Bytes(), cmp)
	}
}
