package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
)

type answer struct {
	Method           string
	URL              *url.URL
	Proto            string // "HTTP/1.0"
	Header           http.Header
	Body             io.ReadCloser
	ContentLength    int64
	TransferEncoding []string
	Host             string
}

func handle(w http.ResponseWriter, r *http.Request) {
	a := answer{
		Method:           r.Method,
		URL:              r.URL,
		Proto:            r.Proto,
		Header:           r.Header,
		Body:             r.Body,
		ContentLength:    r.ContentLength,
		TransferEncoding: r.TransferEncoding,
		Host:             r.Host,
	}
	js, err := json.Marshal(a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	http.HandleFunc("/", handle)
	log.Printf("started on :8080\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
