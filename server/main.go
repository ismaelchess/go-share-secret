package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
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

	setupResponse(&w, r)

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
	setupResponse(&w, r)

	type Result struct {
		Data string
	}

	secret, err := template.ParseFiles("./ui/secret.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	key := mux.Vars(r)["key"]
	if key == "" {
		http.Error(w, "Not complete", http.StatusInternalServerError)
		return
	}

	result, ok := StoreData.Load(key)
	if !ok {
		err1 := secret.Execute(w, &Result{
			Data: "Not exists Data",
		})
		if err1 != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err = secret.Execute(w, &Result{
		Data: result.(string),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	StoreData.Delete(key)

}
func enabledCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func setupResponse(w *http.ResponseWriter, r *http.Request) {

	maxAgeInSeconds := strconv.FormatInt(int64((time.Hour*24)/time.Second), 10)
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
	(*w).Header().Set("Access-Control-Max-Age", maxAgeInSeconds)

}
