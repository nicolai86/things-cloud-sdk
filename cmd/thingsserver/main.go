package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
)

type loggingTransport struct {
	http.RoundTripper
}

func (t *loggingTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	resp, err = t.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}
	body := ioutil.NopCloser(bytes.NewReader(b))
	resp.Body = body
	log.Printf("> %s\n", string(b))
	return resp, nil
}

func main() {
	listen := flag.String("listen", ":443", "port to listen on")
	flag.Parse()

	director := func(req *http.Request) {
		req.URL.Scheme = "https"
		req.Host = "cloud.culturedcode.com"
		req.URL.Host = "cloud.culturedcode.com"
		req.Header.Set("Connection", "close")

		dump, _ := httputil.DumpRequest(req, true)
		log.Printf("< %s\n", string(dump))
	}
	proxy := &httputil.ReverseProxy{Director: director, Transport: &loggingTransport{http.DefaultTransport}}
	log.Printf("Listening on %s\n", *listen)

	err := http.ListenAndServeTLS(*listen, "server.crt", "server.key", proxy)
	// err := http.ListenAndServe(*listen, proxy)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
