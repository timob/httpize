package httpize

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Http Settings
var settings = new(Settings)

// A type that will be used as a httpize.Caller
type CommonFunc func(map[string]Arg) (io.WriterTo, error)

func (f CommonFunc) Call(args map[string]Arg) (io.WriterTo, *Settings, error) {
	writerTo, err := f(args)
	return writerTo, settings, err
}

func Echo(args map[string]Arg) (io.WriterTo, error) {
	name := args["name"].(SafeString)
	return bytes.NewBufferString("Echo " + string(name)), nil
}

var _ = Handle("/Echo?name SafeString", CommonFunc(Echo))

func Greeting(args map[string]Arg) (io.WriterTo, error) {
	return bytes.NewBufferString("Hello World"), nil
}

var _ = Handle("/Greeting", CommonFunc(Greeting))

func ThreeOhThree(args map[string]Arg) (io.WriterTo, error) {
	err := Non500Error{303, "See Other", "http://lookhere"}
	return nil, err
}

var _ = Handle("/ThreeOhThree", CommonFunc(ThreeOhThree))

var count int = 0

func checkCode(t *testing.T, r *httptest.ResponseRecorder, code int) {
	if r.Code != code {
		t.Fatalf("%d %d %v %s", count, r.Code, r.HeaderMap, r.Body)
	}
	t.Logf("%d %d %v %s", count, r.Code, r.HeaderMap, r.Body)
	count++
}

func TestTestApiProvider(t *testing.T) {

	settings.SetToDefault()
	h := GetHandlerForPattern("/Echo?name SafeString")

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "http://host/Echo?name=Gopher", nil)
	h.ServeHTTP(recorder, request)
	checkCode(t, recorder, 200)
	if recorder.Body.String() != "Echo Gopher" {
		t.Fatal("incorrect response")
	}
	if v, ok := recorder.HeaderMap["Content-Type"]; !ok || v[0] != "text/html" {
		t.Fatalf("Content-Type header missing or invalid")
	}

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("GET", "http://host/Echo?badparam=Gopher", nil)
	h.ServeHTTP(recorder, request)
	checkCode(t, recorder, 500)

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("GET", "http://host/Echo?name=Go'pher", nil)
	h.ServeHTTP(recorder, request)
	checkCode(t, recorder, 500)

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("GET", "http://host/Echo?name=Gopher&name=Gopher", nil)
	h.ServeHTTP(recorder, request)
	checkCode(t, recorder, 500)

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("GET", "http://host/Echo?name=Gopher&badparam=Gopher", nil)
	h.ServeHTTP(recorder, request)
	checkCode(t, recorder, 500)

	h = GetHandlerForPattern("/Greeting")

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("GET", "http://host/Greeting", nil)
	h.ServeHTTP(recorder, request)
	checkCode(t, recorder, 200)

	h = GetHandlerForPattern("/ThreeOhThree")

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("GET", "http://host/ThreeOhThree", nil)
	h.ServeHTTP(recorder, request)
	checkCode(t, recorder, 303)

	h = GetHandlerForPattern("/Greeting")

	settings.Cache = 300
	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("GET", "http://host/Greeting", nil)
	h.ServeHTTP(recorder, request)
	checkCode(t, recorder, 200)
	if _, ok := recorder.HeaderMap["Expires"]; !ok {
		t.Fatalf("Expires header missing")
	}
	now := time.Now()
	cacheTime, err := time.Parse(time.RFC1123, recorder.HeaderMap["Expires"][0])
	if err != nil || cacheTime.Before(now) {
		t.Fatalf("Expires header invalid")
	}

	settings.SetToDefault()
	settings.Gzip = true

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("GET", "http://host/Greeting", nil)
	request.Header.Add("Accept-Encoding", "gzip")
	h.ServeHTTP(recorder, request)
	checkCode(t, recorder, 200)
	if v, ok := recorder.HeaderMap["Content-Encoding"]; !ok || v[0] != "gzip" {
		t.Fatalf("Content-Encoding header missing or invalid")
	}

	settings.SetToDefault()

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("GET", "http://host/Greeting", nil)
	h.ServeHTTP(recorder, request)
	checkCode(t, recorder, 200)
	if _, ok := recorder.HeaderMap["Content-Encoding"]; ok {
		t.Fatalf("Unexpected Content-Encoding")
	}
}
