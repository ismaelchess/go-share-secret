package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/ismaelchess/go-share-secret/stores"
	"github.com/ismaelchess/go-share-secret/svc"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	var store stores.Store //= &stores.MapSyncStore{}

	r := mux.NewRouter()

	// Get HOST and PORT
	redisHost := os.Getenv("REDIS_HOST")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	// Create redis connection
	ctx := context.Background()
	store = stores.NewRedisStore(redisHost+":6379", ctx)

	tbl := template.Must(template.ParseFiles("./ui/index.html"))
	parser := &svc.DefaultTemplateParser{}

	r.HandleFunc("/", svc.GoSecret(host, tbl))
	r.HandleFunc("/secret", svc.PostGoSecret(store, host)).Methods(http.MethodPost)
	r.HandleFunc("/secret/{key}", svc.GetGoSecret(store, parser)).Methods(http.MethodGet)
	http.Handle("/", r)

	log.Println("Starting server at port:" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
