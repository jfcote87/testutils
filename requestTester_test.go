// Copyright 2019 James Cote
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testutils

import (
	"net/http"
	"reflect"
	"testing"
)

func TestRequestTester_Check(t *testing.T) {
	type fields struct {
		Path         string
		Auth         string
		Method       string
		Query        string
		Host         string
		ContentType  string
		Header       http.Header
		Payload      []byte
		Response     *http.Response
		ResponseFunc func(*http.Request) (*http.Response, error)
	}
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RequestTester{
				Path:         tt.fields.Path,
				Auth:         tt.fields.Auth,
				Method:       tt.fields.Method,
				Query:        tt.fields.Query,
				Host:         tt.fields.Host,
				ContentType:  tt.fields.ContentType,
				Header:       tt.fields.Header,
				Payload:      tt.fields.Payload,
				Response:     tt.fields.Response,
				ResponseFunc: tt.fields.ResponseFunc,
			}
			if err := r.Check(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("RequestTester.Check() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransport_RoundTrip(t *testing.T) {
	type fields struct {
		Queue []*RequestTester
	}
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *http.Response
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := &Transport{
				Queue: tt.fields.Queue,
			}
			got, err := tx.RoundTrip(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transport.RoundTrip() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Transport.RoundTrip() = %v, want %v", got, tt.want)
			}
		})
	}
}
