package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/ismaelchess/go-share-secret/stores"
	"github.com/ismaelchess/go-share-secret/svc"
	"github.com/subosito/gotenv"
)

func init() {

	if err := gotenv.Load(".env"); err != nil {
		fmt.Println("Error:" + err.Error())
		return
	}
}
func main() {
	//= &MapSyncStore{}
	var store stores.Store

	r := mux.NewRouter()

	// Get HOST and PORT
	redisHost := os.Getenv("REDIS_HOST")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	fmt.Println(redisHost)
	fmt.Println(host)
	fmt.Println(port)

	// // Create redis connection
	ctx := context.Background()
	store = stores.NewRedisStore(redisHost+":6379", ctx)

	tbl := template.Must(template.ParseFiles("./ui/index.html"))
	parser := &svc.DefaultTemplateParser{}

	r.HandleFunc("/", svc.GoSecret(host, tbl))
	r.HandleFunc("/secret", svc.PostGoSecret(store, host)).Methods(http.MethodPost)
	r.HandleFunc("/secret/{key}", svc.GetGoSecret(store, parser)).Methods(http.MethodGet)

	fmt.Println("Starting server at port:" + port)
	panic(http.ListenAndServe(":"+port, r))
}
