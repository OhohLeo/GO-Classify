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
	Method string
	Url    string
	Header map[string]string
	Body   interface{}

	Callback Callback
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

func SendSimple(method string, url string, rsp interface{}) error {

	if pool == nil {
		return fmt.Errorf(
			"%s %s: pool not initialised", method, url)
	}

	return pool.send(&Request{
		Method: method,
		Url:    "http://" + url,
	}, rsp)
}

// Send the requests and wait for the answer
func (p *RequestsPool) send(r *Request, rsp interface{}) (err error) {

	var body io.Reader
	if r.Body != nil {

		var b []byte
		b, err = json.Marshal(r.Body)
		if err != nil {
			return
		}

		body = strings.NewReader(string(b))
	}

	log.Printf("--> %s %s", r.Method, r.Url)
	if p.debug {
	}

	// Send the request
	var req *http.Request
	req, err = http.NewRequest(r.Method, r.Url, body)
	if err != nil {
		return
	}

	// Receive the answer
	httpRsp, err := p.client.Do(req)

	if err != nil {
		log.Printf("--> FAILED %s", err.Error())
		return
	}

	// Get the buffer
	buf := new(bytes.Buffer)
	buf.ReadFrom(httpRsp.Body)

	var debugStr string
	if p.debug {
		debugStr = fmt.Sprintf("Body:%s", buf.String())
	}

	log.Printf("<-- %s%s", httpRsp.Status, debugStr)

	if rsp != nil {
		err = json.Unmarshal(buf.Bytes(), rsp)
	}

	return
}
