package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type Callback interface {
	OnResponse(status int, contentType string, rsp interface{}) error
}

type Request struct {
	Method  string
	Url     string
	Headers map[string]string
	Body    interface{}

	Callback Callback
}

type Response struct {
	Status int
	Body   interface{}
}

type RequestsPool struct {
	client      *http.Client
	requestsMax int
	requests    chan Request
	debug       bool
}

// Uniq instance of the pool of requests
var pool *RequestsPool

// Create the pool of requests with limited size
func New(sizeMax int, debug bool) *RequestsPool {

	pool = &RequestsPool{
		client:      &http.Client{},
		requestsMax: sizeMax,
		debug:       debug,
	}

	return pool
}

func Send(method string, url string, headers map[string]string, body interface{}, rsp interface{}) (chan *Response, error) {

	if pool == nil {
		return nil, fmt.Errorf(
			"%s %s: pool not initialised", method, url)
	}

	return pool.send(&Request{
		Method:  method,
		Url:     url,
		Body:    body,
		Headers: headers,
	}, rsp)
}

// Send the requests and wait for the answer
func (p *RequestsPool) send(r *Request, rsp interface{}) (res chan *Response, err error) {

	var body io.Reader
	if r.Body != nil {

		var b []byte
		b, err = json.Marshal(r.Body)
		if err != nil {
			return
		}

		body = strings.NewReader(string(b))
	}

	var debugStr string
	if p.debug {
		debugStr = fmt.Sprintf(" body: %s", body)
	}
	log.Printf("--> %s %s %s", r.Method, r.Url, debugStr)

	// Create the request
	req, err := http.NewRequest(r.Method, r.Url, body)
	if err != nil {
		return
	}

	// Add headers
	for key, value := range r.Headers {
		req.Header.Add(key, value)
	}

	res = make(chan *Response)

	go func() {

		// Receive the answer
		httpRsp, err := p.client.Do(req)

		if err != nil {
			log.Printf("--> FAILED %s", err.Error())
			close(res)
			return
		}

		// Get the buffer
		buf := new(bytes.Buffer)
		buf.ReadFrom(httpRsp.Body)

		if p.debug {
			debugStr = fmt.Sprintf(" body: %s", buf.String())
		}

		log.Printf("<-- %s%s", httpRsp.Status, debugStr)

		if rsp != nil {
			err = json.Unmarshal(buf.Bytes(), rsp)
		}

		res <- &Response{
			Status: httpRsp.StatusCode,
		}
	}()

	return
}
