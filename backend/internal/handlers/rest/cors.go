package rest

import (
	"net/http"
	"slices"
	"strings"
)

var corsMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
var corsHeaders = []string{"Authorization", "Content-Type"}

func CORS(next http.Handler) http.Handler { return CORSWithOrigins(nil)(next) }

func CORSWithOrigins(origins []string) Middleware {
	allowed := make(map[string]struct{}, len(origins))
	for _, origin := range origins {
		if origin = strings.TrimSpace(origin); origin != "" {
			allowed[origin] = struct{}{}
		}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := allowed[r.Header.Get("Origin")]; ok {
				w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
				w.Header().Add("Vary", "Origin")
			}
			next.ServeHTTP(w, r)
		})
	}
}

func HandleCORSPreFlight(w http.ResponseWriter, r *http.Request, origins []string) {
	if !slices.Contains(origins, r.Header.Get("Origin")) || !slices.Contains(corsMethods, r.Header.Get("Access-Control-Request-Method")) || !requestedHeadersAllowed(r.Header.Get("Access-Control-Request-Headers")) {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Add("Vary", "Origin")
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(corsMethods, ", "))
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(corsHeaders, ", "))
	w.WriteHeader(http.StatusNoContent)
}

func requestedHeadersAllowed(value string) bool {
	for _, header := range strings.Split(value, ",") {
		header = strings.TrimSpace(header)
		if header != "" && !slices.ContainsFunc(corsHeaders, func(allowed string) bool { return strings.EqualFold(header, allowed) }) {
			return false
		}
	}
	return true
}
