package rest

import (
	"net/http"
	"strings"
)

var corsMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
var corsHeaders = []string{"Authorization", "Content-Type"}

func CORS(next http.Handler) http.Handler { return CORSWithOrigins([]string{"*"})(next) }

func CORSWithOrigins(origins []string) Middleware {
	allowed := prepareOriginLookupTable(origins)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := allowed["*"]; ok {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			} else if _, ok := allowed[r.Header.Get("Origin")]; ok {
				w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
				w.Header().Add("Vary", "Origin")
			}

			next.ServeHTTP(w, r)
		})
	}
}

func CORSPreFlight(origins []string) http.HandlerFunc {
	allowed := prepareOriginLookupTable(origins)

	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := allowed[r.Header.Get("Origin")]; ok {
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		}

		w.Header().Add("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(corsMethods, ", "))
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(corsHeaders, ", "))
		w.WriteHeader(http.StatusNoContent)
	}
}

func prepareOriginLookupTable(origins []string) map[string]struct{} {
	allowed := make(map[string]struct{}, len(origins))

	for _, origin := range origins {
		if origin = strings.TrimSpace(origin); origin != "" {
			allowed[origin] = struct{}{}
		}
	}
	return allowed
}
