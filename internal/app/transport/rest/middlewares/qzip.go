package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/kotche/url-shortening-service/internal/app/config"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// GzipHandler decrypts data from the client, encrypts from the server
func GzipHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), config.Compression) {
			reader, err := gzip.NewReader(r.Body)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			r.Body = reader
			defer r.Body.Close()
		}

		if strings.Contains(r.Header.Get("Accept-Encoding"), config.Compression) {
			gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			defer gz.Close()
			w.Header().Set("Content-Encoding", config.Compression)
			w = gzipWriter{ResponseWriter: w, Writer: gz}
		}

		next.ServeHTTP(w, r)
	})
}
