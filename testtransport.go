// Copyright 2017 James Cote and Liberty Fund, Inc.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package testutils contains routines to help with
// creating tests
package testutils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Transport contains an array of http.Response, handler funcs and errors
// that help create an http request test
type Transport struct {
	Queue []interface{}
}

// RoundTrip fulfills the http.Transport interface{} by creating
// http responses and errors from the Transport queue
func (tx *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	var t interface{}
	if len(tx.Queue) > 0 {
		defer func() {
			tx.Queue[0] = nil
			tx.Queue = tx.Queue[1:]
		}()
		t = tx.Queue[0]
	}

	switch t := t.(type) {
	case *http.Response:
		return t, nil
	case func(req *http.Request) (*http.Response, error):
		res, err := t(req)
		return res, err
	case nil:
		return nil, fmt.Errorf("Empty Response")
	default:
		return nil, fmt.Errorf("%v", t)
	}
}

// Add adds a new response to the queue.
func (tx *Transport) Add(val interface{}) {
	tx.Queue = append(tx.Queue, val)
	return
}

// AddResponse creates a response to the queue
func (tx *Transport) AddResponse(status int, body []byte, header http.Header) {
	tx.Add(makeResponse(status, body, header))
}

func makeResponse(status int, body []byte, header http.Header) *http.Response {
	return &http.Response{
		Body:          ioutil.NopCloser(bytes.NewReader(body)),
		StatusCode:    status,
		Status:        http.StatusText(status),
		ContentLength: int64(len(body)),
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        header,
	}
}
