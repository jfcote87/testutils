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

// RequestTester contains expected values for checking against request
type RequestTester struct {
	Path        string
	Auth        string
	Method      string
	Query       string
	Host        string
	ContentType string
	Payload     []byte
	Response    *http.ReadResponse
}

// Check compares expected values with the req parameter
func (r RequestTester) Check(req *http.Request) error {
	if r.Path > "" && r.Path != req.URL.Path {
		return fmt.Errorf("expected request path %s; go %s", r.Path, req.URL.Path)
	}
	if r.Auth > "" && r.Auth != req.Header.Get("Authorization") {
		return fmt.Errorf("expecte auth header %s; got %s", r.Auth, req.Header.Get("Authorization"))
	}
	if r.Method > "" && r.Method != req.Method {
		return fmt.Errorf("expected method %s; got %s", r.Method, req.Method)
	}
	if r.Query > "" && r.Query != req.URL.RawQuery {
		return fmt.Errorf("expected query args %s; got %s", r.Query, req.URL.RawQuery)
	}
	if r.Host > "" && r.Host != req.URL.Host {
		return fmt.Errorf("expected host %s; got %s", r.Host, req.URL.Host)
	}
	if r.ContentType > "" && r.ContentType != req.Header.Get("ContentType") {
		return fmt.Errorf("expected content-type %s; got %s", r.ContentType, req.Header.Get("ContentType"))
	}
	if len(r.Payload) > 0 {
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return fmt.Errorf("unable to read request body: %v", err)
		}
		if bytes.Compare(b, r.Payload) != 0 {
			return fmt.Errorf("expected body %s; got %s", string(r.Payload), string(b))
		}

	}
	return nil
}

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
	case *RequestTester:
		if err := t.Check(req); err != nil {
			return nil, err
		}
		return t.Response, nil
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

// MakeResponse creates an *http.Response for later processing
func MakeResponse(status int, body []byte, header http.Header) *http.Response {
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
