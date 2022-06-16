package middlewares

import (
	"dbms/internal/pkg/utils/log"

	"fmt"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.InfoWithoutCaller(fmt.Sprintf("[%s] %s, %s %s",
			r.Method, r.RemoteAddr, r.URL.Path, time.Since(start)))
	})
}
