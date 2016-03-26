package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var (
	counter int
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

	/*Return a code
	example: ... -H "RETURN_CODE: 503"...
	will respond with status_code of 503
	*/
	if r.Header.Get("RETURN_CODE") != "" {
		header, err := strconv.Atoi(r.Header.Get("RETURN_CODE"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Error(w, "", header)
	}

	/*Return CODE [0] and after [1] tries return [2]
	example: ... -H "RETURN_CODE_AFTER: 500,2,200
	will respond with status_code 500 twice, followed by 200
	*/
	if r.Header.Get("RETURN_CODE_AFTER") != "" {
		vals := strings.Split(r.Header.Get("RETURN_CODE_AFTER"), ",")
		initHeader, err := strconv.Atoi(vals[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tries, err := strconv.Atoi(vals[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		finHeader, err := strconv.Atoi(vals[2])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if counter < tries {
			http.Error(w, "", initHeader)
			counter++
		} else {
			http.Error(w, "", finHeader)
			counter = 0
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	http.HandleFunc("/", handle)
	log.Printf("started on :8080\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
