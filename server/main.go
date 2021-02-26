package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

var StoreData sync.Map

func init() {

}

func main() {

	r := mux.NewRouter()

	fmt.Println("initial")

	tbl := template.Must(template.ParseFiles("./ui/index.html"))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		err := tbl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	})
	r.HandleFunc("/secret", PostGoSecret).Methods(http.MethodPost)
	r.HandleFunc("/secret/{key}", GetGoSecret).Methods(http.MethodGet)

	err := http.ListenAndServe(":8081", r)
	if err != nil {
		log.Fatal(err)
	}
}

func PostGoSecret(w http.ResponseWriter, r *http.Request) {
	var sdata sdata

	if err := json.NewDecoder(r.Body).Decode(&sdata); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	idUrl := sdata.getUniqueId()
	StoreData.Store(idUrl, sdata.Value)

	time.AfterFunc(sdata.expirationDate(), func() {
		StoreData.Delete(idUrl)
	})

	data, err := json.Marshal(&struct {
		URI string `json:"uri"`
	}{
		URI: "http://localhost:8081/secret/" + idUrl,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, _ = w.Write(data)

}

func GetGoSecret(w http.ResponseWriter, r *http.Request) {
	templateExecute := func(w *http.ResponseWriter, m string) error {

		secret, err := template.ParseFiles("./ui/secret.html")
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
		http.Error(w, "Not complete", http.StatusInternalServerError)
		return
	}
	result, ok := StoreData.Load(key)
	if !ok {
		if err := templateExecute(&w, "Not Exist Data"); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if err := templateExecute(&w, result.(string)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	StoreData.Delete(key)
}
