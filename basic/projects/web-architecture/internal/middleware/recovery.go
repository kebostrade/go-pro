package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/DimaJoyti/go-pro/basic/projects/web-architecture/pkg/response"
)

// Recovery is a middleware that recovers from panics
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v\n%s", err, debug.Stack())
				response.InternalServerError(w, "Internal server error")
			}
		}()

		next.ServeHTTP(w, r)
	})
}

