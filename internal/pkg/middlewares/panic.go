package middlewares

import (
	"dbms/internal/pkg/utils/log"

	"fmt"
	"net/http"
)

func PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Warn(fmt.Sprintf("Recovered from panic with err: %s on %s", err, r.RequestURI))
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
