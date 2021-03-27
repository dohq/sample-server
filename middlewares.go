package main

import "net/http"

// BasicAuthMiddleware is apply Basic Auth http Handler.
func (h *Handler) BasicAuthMiddleware(next http.HandlerFunc, username, password string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != username || pass != password {
			w.Header().Set("WWW-Authenticate", "could not exact username and password")
			writeResponse(w, Result{Message: "Not Authorized"}, http.StatusUnauthorized)
			h.Logger.Error("Not Authorized")
			return
		}
		next.ServeHTTP(w, r)
	}
}
