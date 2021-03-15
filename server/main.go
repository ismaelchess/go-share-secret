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
	ctx := context.Background()
	store = NewRedisStore("devredis1:6379", ctx)

	r := mux.NewRouter()
	dataHost := GetPathHost()

	tbl := template.Must(template.ParseFiles("./ui/index.html"))
	parser := &DefaultTemplateParser{}
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		data := &Result{Data: dataHost.Host}
		err := tbl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	r.HandleFunc("/secret", PostGoSecret(store, dataHost)).Methods(http.MethodPost)
	r.HandleFunc("/secret/{key}", GetGoSecret(store, parser)).Methods(http.MethodGet)
	fmt.Println("Starting server at port:" + dataHost.Port)

	panic(http.ListenAndServe(":"+dataHost.Port, r))
}

func PostGoSecret(store Store, pathHost PathHost) http.HandlerFunc {
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
			Path: pathHost.Host + "/" + key,
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

func GetPathHost() PathHost {
	var port = "8081"
	if err := gotenv.Load(); err == nil {
		if port = os.Getenv("PORT"); port == "" {
			port = "8081"
		}
	}
	return PathHost{
		Port: port,
		Host: fmt.Sprintf("http://localhost:%s/secret", port),
	}
}
