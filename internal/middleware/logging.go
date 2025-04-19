package middleware

import (
	"log"
	"net/http"
	"time"
)

// statusRecorder wraps http.ResponseWriter to capture the status code and size.
type statusRecorder struct {
	http.ResponseWriter
	Status int
	Bytes  int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.Status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	if r.Status == 0 {
		// default status if WriteHeader wasn’t called
		r.Status = http.StatusOK
	}
	n, err := r.ResponseWriter.Write(b)
	r.Bytes += n
	return n, err
}

// Logger logs each request’s method, path, remote addr, status, bytes, duration.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w}

		next.ServeHTTP(rec, r)

		log.Printf(
			"%s %s %s %d %d %s %q",
			r.RemoteAddr,
			r.Method,
			r.URL.RequestURI(),
			rec.Status,
			rec.Bytes,
			time.Since(start),
			r.UserAgent(),
		)
	})
}
