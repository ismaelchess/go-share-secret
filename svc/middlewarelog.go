package svc

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func Middlewarelog(sugar *zap.SugaredLogger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		uri := r.RequestURI

		next.ServeHTTP(w, r)
		duration := time.Since(start)

		sugar.Infow("Method URL :",
			"url", uri,
			"Method", r.Method,
			"duration", duration,
		)

	})
}
