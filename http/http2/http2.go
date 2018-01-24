package http2

import (
	"log"
	"net/http"
)

// New http2 pusher handler.
func New(files []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if p, ok := w.(http.Pusher); ok {
			for i := range files {
				if err := p.Push(files[i], nil); err != nil {
					log.Printf("Failed to push %v: %v", files[i], err)
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}
