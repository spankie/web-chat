package web

import (
	"log"
	"net/http"
	"time"
)

//LoggingHandler Logs request time, method and duration of handler/request execution
func LoggingHandler(next http.Handler) http.Handler {
	logger := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		responseTime := time.Since(t1)
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), responseTime)
	}

	return http.HandlerFunc(logger)
}

func Login(w http.ResponseWriter, r *http.Request) {
	// get the post parameters ...

}
