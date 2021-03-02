package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
)

var StoreData sync.Map
var DataHost PathHost

func init() {
	gotenv.Load()
}

func main() {
	r := mux.NewRouter()
	DataHost = GetPathHost()

	tbl := template.Must(template.ParseFiles("./ui/index.html"))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		data := &Result{Data: DataHost.Host}
		err := tbl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	r.HandleFunc("/secret", PostGoSecret).Methods(http.MethodPost)
	r.HandleFunc("/secret/{key}", GetGoSecret).Methods(http.MethodGet)

	fmt.Println("Run Application, Port:" + DataHost.Port)

	err := http.ListenAndServe(":"+DataHost.Port, r)
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
		Path string `json:"uri"`
	}{
		Path: DataHost.Host + "/" + idUrl,
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

func GetPathHost() PathHost {
	var port string

	if port = os.Getenv("PORT"); port == "" {
		port = "8081"
	}

	return PathHost{
		Port: port,
		Host: fmt.Sprintf("http://localhost:%s/secret", port),
	}
}
