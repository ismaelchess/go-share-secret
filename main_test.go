package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/ismaelchess/go-share-secret/stores"
	"github.com/ismaelchess/go-share-secret/svc"
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
			var store stores.Store = &stores.MapSyncStore{}

			//gotenv.Load()
			host_get := os.Getenv("HOST")

			req, err := http.NewRequest(http.MethodPost, "localhost:8080/secret", strings.NewReader(tc.value))
			if err != nil {
				t.Fatalf("could not created request: %v", err)
			}

			rec := httptest.NewRecorder()
			httpGoSecret := svc.PostGoSecret(store, host_get)
			httpGoSecret.ServeHTTP(rec, req)

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
		{name: "correct value", value: "9f43506d-c5f0-45e9-8871-d513e8da9018", status: http.StatusNotFound, wantErr: true},
		{name: "missing value", value: "", status: http.StatusNotFound, wantErr: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var store stores.Store = &stores.MapSyncStore{}

			req, err := http.NewRequest("GET", "localhost:8080/secret/1", nil)
			if err != nil {
				t.Fatalf("could not created request: %v", err)
			}

			vars := map[string]string{"key": tc.value}
			req = mux.SetURLVars(req, vars)

			rec := httptest.NewRecorder()

			parser := &svc.TestTemplateParser{}
			httpGoSecret := svc.GetGoSecret(store, parser)
			httpGoSecret.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if tc.wantErr {
				if res.StatusCode != http.StatusNotFound {
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
