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
			err:  fmt.Errorf("error"),
			code: http.StatusInternalServerError,
			want: "Got error with details: error",
		},
		{
			name: "error with code 400",
			err:  fmt.Errorf("error"),
			code: http.StatusBadRequest,
			want: "Got error with details: error",
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
			name: "success",
			resp: &bqrfgemini.BigQueryResponse{
				Replies: []string{"success"},
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
			if reflect.DeepEqual(&got.Replies, test.want) {
				t.Errorf("SendSuccess(%v) = %v, want %v", test.resp, got.Replies, test.want)
			}
		})
	}
}
