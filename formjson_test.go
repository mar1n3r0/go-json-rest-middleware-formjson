package formJson

import (
	"bytes"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ant0ine/go-json-rest/rest/test"
)

type JSON map[string]interface{}

func simpleGetEndpoint(w rest.ResponseWriter, r *rest.Request) {
	body := map[string]interface{}{
		"status": "success",
	}
	w.WriteJson(body)
}

func simplePostEndpoint(w rest.ResponseWriter, r *rest.Request) {
	var body map[string]interface{}
	r.DecodeJsonPayload(&body)
	w.WriteJson(body)
}

func NewSimpleAPI(mw *MiddleWare) http.Handler {
	api := rest.NewApi()
	api.Use(mw)
	router, _ := rest.MakeRouter(
		rest.Post("/", simplePostEndpoint),
		rest.Get("/", simpleGetEndpoint),
	)
	api.SetApp(router)
	return api.MakeHandler()
}

func TestPostValidFormData(t *testing.T) {
	handler := NewSimpleAPI(&MiddleWare{})

	data := url.Values{}
	data.Set("name", "foo")
	data.Add("surname", "bar")

	r, _ := http.NewRequest("POST", "http://localhost/", bytes.NewBufferString(data.Encode())) // <-- URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	recordedPost := test.RunRequest(t, handler, r)

	recordedPost.CodeIs(http.StatusOK)
	recordedPost.BodyIs(`{"name":"foo","surname":"bar"}`)
}

func TestPostInvalidFormData(t *testing.T) {
	handler := NewSimpleAPI(&MiddleWare{})

	data := "mal&formed"

	r, _ := http.NewRequest("POST", "http://localhost/", bytes.NewBufferString(data)) // <-- URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data)))

	recordedPost := test.RunRequest(t, handler, r)

	recordedPost.CodeIs(http.StatusInternalServerError)
	recordedPost.BodyIs(`{"Error":"Error Converting Form Data"}`)
}

func TestGet(t *testing.T) {
	handler := NewSimpleAPI(&MiddleWare{})

	req := test.MakeSimpleRequest("GET", "http://localhost/", nil)
	recorded := test.RunRequest(t, handler, req)

	recorded.CodeIs(http.StatusOK)
	recorded.BodyIs(`{"status":"success"}`)
}

func TestPostJSON(t *testing.T) {
	handler := NewSimpleAPI(&MiddleWare{})

	simplePostData := map[string]interface{}{
		"john": "doe",
		"age":  30,
	}

	postRequest := test.MakeSimpleRequest("POST", "https://localhost/", simplePostData)
	recordedPost := test.RunRequest(t, handler, postRequest)

	recordedPost.CodeIs(http.StatusOK)
	recordedPost.BodyIs(`{"age":30,"john":"doe"}`)
}
