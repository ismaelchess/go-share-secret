package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"text/template"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
)

var StoreData sync.Map
var DataHost PathHost

var (
	ctx = context.Background()
	dbr *redis.Client
)

func init() {

	dbr = redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     "devredis1:6379",
		Password: "",
		DB:       0,
	})

	fmt.Println("Connect to DB:", dbr)
}

func main() {
	var store Store

	// quiero probar con mapas
	//store=NewMapStore()
	// ya no quiero mapas, ahora quiero que funcione aunque reinicie mi server
	// voy a probar con redis
	store = NewRedisStore("localhost", "6339")

	r := mux.NewRouter()
	DataHost = GetPathHost()

	tbl := template.Must(template.ParseFiles("./ui/index.html"))
	parser := &DefaultTemplateParser{}
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		data := &Result{Data: DataHost.Host}
		err := tbl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	r.HandleFunc("/secret", PostGoSecret(store)).Methods(http.MethodPost)
	r.HandleFunc("/secret/{key}", GetGoSecret(parser)).Methods(http.MethodGet)

	fmt.Println("Starting server at port:" + DataHost.Port)

	panic(http.ListenAndServe(":"+DataHost.Port, r))
}

func PostGoSecret(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sdata sdata

		if err := json.NewDecoder(r.Body).Decode(&sdata); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		idUrl := sdata.getUniqueId()

		err := store.Save(r.Context(), idUrl, sdata.Value, sdata.expirationDate())

		//err := dbr.Set(ctx, idUrl, sdata.Value, sdata.expirationDate()).Err()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

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
}

func GetGoSecret(parser TemplateParser) http.HandlerFunc {
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

		result, err := dbr.Get(ctx, key).Result()
		if err == redis.Nil {
			if err := templateExecute(&w, "Not Exist Data"); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := templateExecute(&w, result); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		dbr.Del(ctx, key)
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
