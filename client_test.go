package thingscloud

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
)

func stringVal(str string) *string {
	return &str
}

type fakeResponse struct {
	statusCode int
	file       string
}

func fakeServer(t fakeResponse) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Open(fmt.Sprintf("tapes/%s", t.file))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Printf("Unable to open fixture %q path %q", t.file, r.URL.Path)
			return
		}
		defer f.Close()

		content, err := ioutil.ReadAll(f)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Printf("Unable to load fixture %q path %q", t.file, r.URL.Path)
			return
		}
		w.WriteHeader(t.statusCode)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, string(content))
	}))
}
