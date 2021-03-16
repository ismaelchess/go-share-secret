package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
)

func main() {
	var store Store //= &MapSyncStore{}

	gotenv.Load()
	r := mux.NewRouter()

	// Get HOST and PORT
	redis_host := os.Getenv("REDIS_HOST")
	redis_port := os.Getenv("REDIS_PORT")
	host_get := os.Getenv("HOST_GET")
	port := os.Getenv("PORT")

	// Create redis connection
	ctx := context.Background()
	store = NewRedisStore(redis_host+":"+redis_port, ctx)

	tbl := template.Must(template.ParseFiles("./ui/index.html"))
	parser := &DefaultTemplateParser{}
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		data := &Result{Data: host_get}
		err := tbl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	r.HandleFunc("/secret", PostGoSecret(store, host_get)).Methods(http.MethodPost)
	r.HandleFunc("/secret/{key}", GetGoSecret(store, parser)).Methods(http.MethodGet)

	fmt.Println("Starting server at port:" + port)
	panic(http.ListenAndServe(":"+port, r))
}

func PostGoSecret(store Store, hostGet string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sdata sdata

		if err := json.NewDecoder(r.Body).Decode(&sdata); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		key := sdata.getUniqueId()
		err := store.Save(key, sdata.Value, sdata.expirationDate())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		data, err := json.Marshal(&struct {
			Path string `json:"uri"`
		}{
			Path: hostGet + "/" + key,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, _ = w.Write(data)
	}
}

func GetGoSecret(store Store, parser TemplateParser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templateExecute := func(w *http.ResponseWriter, m string) error {
			secret, err := parser.ParseFiles()
			if err != nil {
				return err
			}
			err1 := secret.Execute(*w, &Result{
				Data: m,
			})
			if err1 != nil {
				return err
			}
			return nil
		}
		key := mux.Vars(r)["key"]
		if key == "" {
			http.Error(w, "Not complete", http.StatusBadRequest)
			return
		}
		value, err := store.Load(key)
		if value == "" && err == nil {
			if err := templateExecute(&w, "Not Exist Data"); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := templateExecute(&w, value); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = store.Delete(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
