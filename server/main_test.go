package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestPostGoSecret(t *testing.T) {
	tt := []struct {
		name    string
		value   string
		status  int
		wantErr bool
	}{
		{name: "", value: `{"value": "hola", "unit": "m", "utime":2}`, status: http.StatusOK, wantErr: false},
		{name: "missing value", value: "", status: http.StatusBadRequest, wantErr: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			var store Store = &MapSyncStore{}
			dataHost := GetPathHost()

			req, err := http.NewRequest(http.MethodPost, "localhost:8081/secret", strings.NewReader(tc.value))
			if err != nil {
				t.Fatalf("could not created request: %v", err)
			}

			rec := httptest.NewRecorder()
			httpGoSecret := PostGoSecret(store, dataHost)
			httpGoSecret(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if tc.wantErr {
				if res.StatusCode != http.StatusBadRequest {
					t.Errorf("expected status badrequest; got %v", res.Status)
				}
			} else {
				if res.StatusCode != http.StatusOK {
					t.Errorf("expected status ok; got %v", res.Status)
				}
				b, err := ioutil.ReadAll(res.Body)
				if err != nil {
					t.Fatalf("could not read response: %v", err)
				}
				r := string(bytes.TrimSpace(b))
				if r == "" {
					t.Fatalf("there is no value: %s", r)
				}
			}

		})

	}
}

func TestGetGoSecret(t *testing.T) {
	tt := []struct {
		name    string
		value   string
		status  int
		wantErr bool
	}{
		{name: "correct value", value: "9f43506d-c5f0-45e9-8871-d513e8da9018", status: http.StatusOK, wantErr: false},
		{name: "missing value", value: "", status: http.StatusBadRequest, wantErr: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var store Store = &MapSyncStore{}

			req, err := http.NewRequest("GET", "localhost:8081/secret", nil)
			if err != nil {
				t.Fatalf("could not created request: %v", err)
			}

			vars := map[string]string{"key": tc.value}
			req = mux.SetURLVars(req, vars)

			rec := httptest.NewRecorder()

			parser := &TestTemplateParser{}
			httpGoSecret := GetGoSecret(store, parser)
			httpGoSecret(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if tc.wantErr {
				if res.StatusCode != http.StatusBadRequest {
					t.Errorf("expected status badrequest; got %v", res.Status)
				}
			} else {
				if res.StatusCode != http.StatusOK {
					t.Errorf("expected status ok; got %v", res.Status)
				}
				b, err := ioutil.ReadAll(res.Body)
				if err != nil {
					t.Fatalf("could not read response: %v", err)
				}
				r := string(bytes.TrimSpace(b))
				if r == "" {
					t.Fatalf("there is no value: %s", r)
				}
			}
		})
	}
}
