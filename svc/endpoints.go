package svc

import (
	"encoding/json"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/ismaelchess/go-share-secret/stores"
)

func GoSecret(host string, tbl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := &Result{Data: host + "/secret"}
		err := tbl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func PostGoSecret(store stores.Store, host string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var par Par

		if err := json.NewDecoder(r.Body).Decode(&par); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		key := par.getUniqueId()
		err := store.Save(key, par.Value, par.expirationDate())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		data, err := json.Marshal(&struct {
			Path string `json:"uri"`
		}{
			Path: host + "/secret/" + key,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = w.Write(data)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	})
}

func GetGoSecret(store stores.Store, parser TemplateParser) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		executeTemplate := func(w *http.ResponseWriter, m string) error {
			secret, err := parser.ParseFiles()
			if err != nil {
				return err
			}
			err = secret.Execute(*w, &Result{
				Data: m,
			})
			if err != nil {
				return err
			}
			return nil
		}
		key := mux.Vars(r)["key"]
		value, err := store.Load(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err = executeTemplate(&w, value); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err = store.Delete(key); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
