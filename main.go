package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/ismaelchess/go-share-secret/stores"
	"github.com/ismaelchess/go-share-secret/svc"
	zap "go.uber.org/zap"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	var store stores.Store //= &stores.MapSyncStore{}

	logger := zap.NewExample()
	//defer logger.Sync()
	defer func() {
		err := logger.Sync()
		if err != nil {
			log.Panic(err)
		}
	}()
	sugar := logger.Sugar()

	mux := http.NewServeMux()

	// Get HOST and PORT
	redisHost := os.Getenv("REDIS_HOST")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	// Create redis connection
	ctx := context.Background()
	store = stores.NewRedisStore(redisHost+":6379", ctx)

	tbl := template.Must(template.ParseFiles("./ui/index.html"))
	parser := &svc.DefaultTemplateParser{}

	mux.Handle("/", svc.Middlewarelog(sugar, svc.GoSecret(host, tbl)))
	mux.Handle("/secret", svc.Middlewarelog(sugar, svc.PostGoSecret(store, host)))        //.Methods(http.MethodPost)
	mux.Handle("/secret/{key}", svc.Middlewarelog(sugar, svc.GetGoSecret(store, parser))) //.Methods(http.MethodGet)

	sugar.Info("Initial Server")

	log.Println("Starting server at port:" + port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
