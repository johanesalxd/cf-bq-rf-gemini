package bqrfgemini_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	bqrfgemini "github.com/johanesalxd/cf-bq-rf-gemini"
)

func TestSendError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		code int
		want string
	}{
		{
			name: "error with code 500",
			err:  fmt.Errorf("internal server error"),
			code: http.StatusInternalServerError,
			want: "Got error with details: internal server error",
		},
		{
			name: "error with code 400",
			err:  fmt.Errorf("bad request"),
			code: http.StatusBadRequest,
			want: "Got error with details: bad request",
		},
		{
			name: "error with code 404",
			err:  fmt.Errorf("not found"),
			code: http.StatusNotFound,
			want: "Got error with details: not found",
		},
		{
			name: "error with code 403",
			err:  fmt.Errorf("forbidden"),
			code: http.StatusForbidden,
			want: "Got error with details: forbidden",
		},
		{
			name: "error with code 401",
			err:  fmt.Errorf("unauthorized"),
			code: http.StatusUnauthorized,
			want: "Got error with details: unauthorized",
		},
		{
			name: "error with custom message",
			err:  fmt.Errorf("custom error message"),
			code: http.StatusInternalServerError,
			want: "Got error with details: custom error message",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			bqrfgemini.SendError(w, test.err, test.code)

			resp := w.Result()
			if resp.StatusCode != test.code {
				t.Errorf("SendError(%v, %v) = %v, want %v", test.err, test.code, resp.StatusCode, test.code)
			}
			if resp.Header.Get("Content-Type") != "application/json" {
				t.Errorf("SendError(%v, %v) = %v, want %v", test.err, test.code, resp.Header.Get("Content-Type"), "application/json")
			}

			var got bqrfgemini.BigQueryResponse
			if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
				t.Errorf("SendError(%v, %v) = %v, want %v", test.err, test.code, err, test.want)
			}
			if got.ErrorMessage != test.want {
				t.Errorf("SendError(%v, %v) = %v, want %v", test.err, test.code, got.ErrorMessage, test.want)
			}
		})
	}
}

func TestSendSuccess(t *testing.T) {
	tests := []struct {
		name string
		resp *bqrfgemini.BigQueryResponse
		want []string
	}{
		{
			name: "single success reply",
			resp: &bqrfgemini.BigQueryResponse{
				Replies: []string{"success"},
			},
			want: []string{"success"},
		},
		{
			name: "multiple success replies",
			resp: &bqrfgemini.BigQueryResponse{
				Replies: []string{"success1", "success2", "success3"},
			},
			want: []string{"success1", "success2", "success3"},
		},
		{
			name: "empty replies",
			resp: &bqrfgemini.BigQueryResponse{
				Replies: []string{},
			},
			want: []string{},
		},
		{
			name: "nil replies",
			resp: &bqrfgemini.BigQueryResponse{
				Replies: nil,
			},
			want: nil,
		},
		{
			name: "response with error message",
			resp: &bqrfgemini.BigQueryResponse{
				Replies:      []string{"success"},
				ErrorMessage: "Some error occurred",
			},
			want: []string{"success"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			bqrfgemini.SendSuccess(w, test.resp)

			resp := w.Result()
			if resp.StatusCode != http.StatusOK {
				t.Errorf("SendSuccess(%v) = %v, want %v", test.resp, resp.StatusCode, http.StatusOK)
			}
			if resp.Header.Get("Content-Type") != "application/json" {
				t.Errorf("SendSuccess(%v) = %v, want %v", test.resp, resp.Header.Get("Content-Type"), "application/json")
			}

			var got bqrfgemini.BigQueryResponse
			if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
				t.Errorf("SendSuccess(%v) = %v, want %v", test.resp, err, test.want)
			}
			if !reflect.DeepEqual(got.Replies, test.want) {
				t.Errorf("SendSuccess(%v) = %v, want %v", test.resp, got.Replies, test.want)
			}
			if got.ErrorMessage != test.resp.ErrorMessage {
				t.Errorf("SendSuccess(%v) ErrorMessage = %v, want %v", test.resp, got.ErrorMessage, test.resp.ErrorMessage)
			}
		})
	}
}
