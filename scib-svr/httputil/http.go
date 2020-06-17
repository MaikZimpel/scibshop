package httputil

import (
	"net/http"
	"time"
)

// addCookie will apply a new cookie to the response of a http request
// with the key/value specified.
func addCookie(w http.ResponseWriter, name, value string, ttl time.Duration) {
	expire := time.Now().Add(ttl)
	cookie := http.Cookie{
		Name:    name,
		Value:   value,
		Expires: expire,
	}
	http.SetCookie(w, &cookie)
}

func EnableCors(w *http.ResponseWriter, address string) {
	(*w).Header().Set("Access-Control-Allow-Origin", address)
}
