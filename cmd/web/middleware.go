package main

import (

	"net/http"

)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Content-Security-Policy",
			       "default-src 'self'; "+
			       "style-src 'self' 'unsafe-inline'; "+
			       "script-src 'self' 'unsafe-inline' https://unpkg.com")

		next.ServeHTTP(w,r)
	})
}


func blockStaticFolder(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		if r.URL.Path == "/static/" {
			http.Error(w, "Access Not Allowed...", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
