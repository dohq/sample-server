package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"go.uber.org/zap"
)

var h = NewHandler(os.Stdout, zap.ErrorLevel)

func TestHandler_GetStatus(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h.GetStatus(w, r)

	result, err := unmarshaler(w.Body)
	if err != nil {
		t.Errorf("could not unmarshal response body: %v", err.Error())
	}

	if w.Code != http.StatusOK {
		t.Errorf("got %v, want %v", w.Code, http.StatusOK)
	}

	if result.Message != "ok" {
		t.Errorf("got %v, want %v", result.Message, "ok")
	}
}

func TestHandler_GetTime(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h.GetTime(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("got %v, want %v", w.Code, http.StatusOK)
	}
}

func TestHandler_GetSleep(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		in           string
		responseBody string
		statusCode   int
		duration     float64
	}{
		{
			name:         "success: normal",
			in:           "1s",
			responseBody: "1s",
			statusCode:   http.StatusOK,
			duration:     1,
		},
		{
			name:         "success: duration zero",
			in:           "0s",
			responseBody: "0s",
			statusCode:   http.StatusOK,
			duration:     0,
		},
		{
			name:         "fail: null",
			in:           "null",
			responseBody: "could not parse query",
			statusCode:   http.StatusInternalServerError,
			duration:     0,
		},
		{
			name:         "fail: blank duration",
			in:           "",
			responseBody: "duration is blank",
			statusCode:   http.StatusInternalServerError,
			duration:     0,
		},
		{
			name:         "fail: parse query",
			in:           "dummy",
			responseBody: "could not parse query",
			statusCode:   http.StatusInternalServerError,
			duration:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/sleep?d="+tt.in, nil)
			w := httptest.NewRecorder()

			start := time.Now()
			h.GetSleep(w, r)
			end := time.Since(start).Seconds()

			if tt.statusCode != w.Code {
				t.Errorf("want %v, got %v", tt.statusCode, w.Code)
			}

			result, err := unmarshaler(w.Body)
			if err != nil {
				t.Errorf("could not unmarshal response body: %v", err.Error())
			}

			if tt.responseBody != result.Message {
				t.Errorf("want %v, got %v", tt.responseBody, result.Message)
			}

			if tt.statusCode == http.StatusOK && 1 <= tt.duration {
				if end < tt.duration || end > tt.duration*1.05 {
					t.Errorf("want %v, got %v", tt.duration, end)
				}
			}
		})
	}
}

// func TestHandler_GetRemoteIP(t *testing.T) {
// 	type fields struct {
// 		Logger *zap.Logger
// 	}
// 	type args struct {
// 		w http.ResponseWriter
// 		r *http.Request
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			h := &Handler{
// 				Logger: tt.fields.Logger,
// 			}
// 			h.GetRemoteIP(tt.args.w, tt.args.r)
// 		})
// 	}
// }
//
// func TestHandler_GetEnv(t *testing.T) {
// 	type fields struct {
// 		Logger *zap.Logger
// 	}
// 	type args struct {
// 		w http.ResponseWriter
// 		r *http.Request
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			h := &Handler{
// 				Logger: tt.fields.Logger,
// 			}
// 			h.GetEnv(tt.args.w, tt.args.r)
// 		})
// 	}
// }

func unmarshaler(buf *bytes.Buffer) (Result, error) {
	var result Result

	b, err := ioutil.ReadAll(buf)
	if err != nil {
		return result, err
	}
	if err := json.Unmarshal(b, &result); err != nil {
		return result, err
	}

	return result, nil
}
