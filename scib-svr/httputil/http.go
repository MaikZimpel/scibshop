package httputil

import (
	"fmt"
	"net/http"
	"strings"
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

// formatRequest generates ascii representation of a request
func FormatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		_ = r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	if r.Method == "PUT" {
		request = append(request, "Body: \n")
		var b []byte
		_, err := r.Body.Read(b)
		if err != nil {
			panic(fmt.Errorf("err parsing body"))
		}
		request = append(request, "\n")
		request = append(request, string(b))
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}
