package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func WithCompressionResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(rw, r)
			return
		}

		gz, gzipError := gzip.NewWriterLevel(rw, gzip.BestSpeed)
		if gzipError != nil {
			http.Error(rw, gzipError.Error(), http.StatusBadRequest)
			return
		}
		defer gz.Close()

		rw.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: rw, Writer: gz}, r)

	})
}

func Unzip(res http.ResponseWriter, req *http.Request) (io.Reader, error) {
	var reader io.Reader

	if req.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(req.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return nil, err
		}
		reader = gz
		defer gz.Close()
	} else {
		reader = req.Body
	}

	return reader, nil
}
