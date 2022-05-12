package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/kotche/url-shortening-service/internal/config"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func GzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if strings.Contains(request.Header.Get("Content-Encoding"), config.Compression) {
			reader, err := gzip.NewReader(request.Body)

			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}

			request.Body = reader
			defer request.Body.Close()
		}

		if strings.Contains(request.Header.Get("Accept-Encoding"), config.Compression) {
			gz, err := gzip.NewWriterLevel(writer, gzip.BestSpeed)

			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}

			defer gz.Close()
			writer.Header().Set("Content-Encoding", config.Compression)
			writer = gzipWriter{ResponseWriter: writer, Writer: gz}
		}

		next.ServeHTTP(writer, request)
	})
}
